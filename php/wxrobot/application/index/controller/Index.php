<?php

namespace app\index\controller;

use app\index\model\Params;
use app\index\model\Synckey;
use app\index\model\User;
use think\Cookie;
use think\Request;
use think\Controller;
use WxApi\WxApi;

class Index extends Controller
{


    /**
     * 请求微信包体
     * @var array BaseRequest
     */
    private $_BaseRequest;
    /**
     * 请求微信uin
     * @var string Uin
     */
    private $_Uin;
    /**
     * 请求微信Sid
     * @var string Sid
     */
    private $_Sid;
    /**
     * 请求微信Skey
     * @var string Skey
     */
    private $_Skey;
    /**
     * 时间+e;
     * @var string DeviceId
     */
    private $_DeviceId;
    /**
     * 心跳检测数据包
     * @var array Synckey
     */
    private $_Synckey;
    /**
     * 实例化接口
     * @var WxApi
     */
    private $WxApi;
    /**
     * @var mixed
     */
    private $_post;
    /**
     * @var string
     */
    private $webwx_data_ticket;
    /**
     * @var string
     */
    private $webwxuvid;
    /**
     * @var string
     */
    private $webwx_auth_ticket;
    /**
     * @var string|Request
     */
    private $action_name;

    /**
     * Index constructor.
     */
    public function __construct()
    {
        $request = Request::instance();
        $this->action_name = $request->action();
        $model = new Params();
        /**
         * 获取存储的回调参数
         */
        $post_data = $model->order('id desc')->find();
        $this->_post = $post_data['params'];
        /**
         * 实例化接口函数
         */
        $this->WxApi = new WxApi($this->_post);
        /**
         * 构造cookie
         */
        $cookie = file(EXTEND_PATH . '/WxApi/wx.cookie');
        $webwx_data_ticket = explode('webwx_data_ticket', $cookie[8]);
        $this->webwx_data_ticket = trim($webwx_data_ticket[1]);
        $webwxuvid = explode('webwxuvid', $cookie[9]);
        $this->webwxuvid = trim($webwxuvid[1]);
        $webwx_auth_ticket = explode('webwx_auth_ticket', $cookie[10]);
        $this->webwx_auth_ticket = trim($webwx_auth_ticket[1]);
        Cookie::set('webwx_data_ticket', $this->webwx_data_ticket, 3600);
        Cookie::set('webwx_auth_ticket', $this->webwx_auth_ticket, 3600);
        Cookie::set('webwxuvid', $this->webwxuvid, 3600);
    }

    public function index()
    {
        $data = $this->login();
        $callback_url = Url('index/index/callback', array('uuid' => $data['uuid']));
        return view('index', ['data' => $data, 'callback_url' => $callback_url]);
    }

    /**
     * 登录
     */
    public function login()
    {
        $data = $this->WxApi->login();
        return $data;
    }

    /**
     * 检查登录
     */
    public function check_login()
    {
        $uuid = input('param.uuid');
        $time = input('param.time');
        if (empty($time)) exit();
        set_time_limit(0);//无限请求超时时间
        $i = 0;
        while (true) {
            //sleep(1);
            usleep(500000);//0.5秒
            $i++;

            //若得到数据则马上返回数据给客服端，并结束本次请求
            $rand = rand(1, 999);
            if ($rand <= 150) {
                echo $this->WxApi->checkLogin($uuid);
                exit();
            }

            //服务器($_POST['time']*0.5)秒后告诉客服端无数据
            if ($i == $time) {
                $arr = array('success' => "0");
                echo json_encode($arr);
                exit();
            }
        }
    }

    /**
     * 回调
     */
    public function callback()
    {
        $uuid = input('param.uuid');

        $callback = $this->WxApi->get_uri($uuid);

        $post = $this->WxApi->post_self($callback);
        if ($post) {

            $model = new Params();

            $model->params = json_encode($post);

            $model->save();
            return $this->success('成功', 'index/wxinit');
        } else {
            return $this->error('失败', 'index/index');

        }
    }

    /**
     *微信初始化
     */
    public function wxinit()
    {

        $data = $this->WxApi->wxinit($this->_post);
        $_ContactList = $data['ContactList'];
        foreach ($_ContactList as $key => $value) {
            $res = preg_match('/@@/', $value['UserName']);
            /**
             * 如果是群组
             */
            if ($res) {
                $_ContactList[$key]['user_type'] = '微信群';
            }else{
                $_ContactList[$key]['user_type'] = '好友';
            }
        }
        /**
         * 存储登录人信息
         */
        if ($data['User']) {
            $model = new User();
            $model->postUser($data['User']['NickName'], $data['User']['UserName'], $data['User']['HeadImgUrl']);
        }
        /**
         * 存储心跳包
         */
        $model = new Synckey();
        $sv = $this->WxApi->ImplodeSynckey($data['SyncKey']['List']);
        $model->compare($data['SyncKey'], $sv, $type = 0);

        $this->WxApi->wxstatusnotify($this->_post);

        $this->WxApi->synccheck($data['SyncKey'], 10);

        if ($data['BaseResponse']['Ret'] === 0) {
            return view('member_list', [
                    'user' => $data['User'],
                    'ContactList' => $_ContactList,
                    'action_name' => $this->action_name,
                    'list' => [],
                ]
            );
        } else {
            return $this->error('失败', 'index/index');
        }
    }

    /**
     * 心跳检测
     */
    public function syncheck()
    {
        $time = input('param.time');
        if (empty($time)) exit();
        set_time_limit(0);//无限请求超时时间
        $i = 0;
        while (true) {
            //sleep(1);
            usleep(500000);//0.5秒
            $i++;

            //若得到数据则马上返回数据给客服端，并结束本次请求
            $rand = rand(1, 999);
            if ($rand <= 150) {
                $model = new Synckey();
                $synckey_data = $model->order('id', 'desc')->find();
                $synckey = $synckey_data['synckey'];
                $synckey = json_decode($synckey, true);
                $data = $this->WxApi->synccheck($synckey, 10);
                $data['success'] = 1;
                echo json_encode($data);
                exit();
            }

            //服务器($_POST['time']*0.5)秒后告诉客服端无数据
            if ($i == $time) {
                $arr = array('success' => "1", 'ret' => 0, 'sel' => 0);
                echo json_encode($arr);
                exit();
            }
        }
    }

    /**
     * 获取新消息
     */
    public function webwxsync()
    {
        $model = new Synckey();
        $synckey_data = $model->order('id', 'desc')->find();
        $synckey = json_decode($synckey_data['synckey'], true);
        $data = $this->WxApi->webwxsync($synckey);
        if (count($data['SyncCheckKey']['List']) > 1) {
            $new_Synckey = $data['SyncCheckKey'];
            $sv = $this->WxApi->ImplodeSynckey($new_Synckey['List']);
            $model->compare($new_Synckey, $sv, 1);
        }
        echo json_encode($data['AddMsgList']);

    }

    public function ContactList()
    {
        return view('member_list', [
                'action_name' => $this->action_name,
                'list' => [],
                'pages' => '',
                'ContactList' => []
            ]
        );
    }

    public function keywords(){

    }
}
