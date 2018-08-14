<?php

namespace WxApi;

use app\index\model\Synckey;

/**
 * Created by PhpStorm.
 * User: marcelo
 * Date: 2017/4/12
 * Time: 11:10
 */
class WxApi
{
    /**
     * 微信二维码ID
     * @var string Uuid
     */
    private $_Uuid;
    /**
     * 微信二维码
     * @var string
     */
    private $_qrcode;
    /**
     * 微信网页版appid
     * @var string
     */
    private $_appid = 'wx782c26e4c19acffb';
    private $_response;
    private $_r;

    function __construct($response)
    {
        $this->_response = json_decode($response, true);

    }

    /**
     * 毫秒级时间戳
     * @return string
     */
    function getMillisecond()
    {
        list($t1, $t2) = explode(' ', microtime());
        return $t2 . ceil(($t1 * 1000));
    }


    /**
     * 获取唯一的uuid用于生成二维码
     * @param $appid
     * @return mixed
     */
    function get_uuid($appid)
    {
        $url = 'https://login.weixin.qq.com/jslogin';
        $url .= '?appid=' . $appid;
        $url .= '&fun=new';
        $url .= '&lang=zh_CN';
        $url .= '&_=' . time();

        $content = $this->curlPost($url);
        //也可以使用正则匹配
        $content = explode(';', $content);

        $content_uuid = explode('"', $content[1]);

        $uuid = $content_uuid[1];

        return $uuid;
    }

    /**
     * 生成二维码
     * @param $uuid
     * @return img
     */
    function get_qrcode($uuid)
    {
        $url = 'https://login.weixin.qq.com/qrcode/' . $uuid . '?t=webwx';
        $img = "<img class='img' src=" . $url . "/>";
        return $img;
    }

    /**
     * 输出登录二维码
     * @return array
     */
    function login()
    {
        $this->_Uuid = $this->get_uuid($this->_appid);
        $this->_qrcode = $this->get_qrcode($this->_Uuid);
        return [
            'uuid' => $this->_Uuid,
            'qrcode' => $this->_qrcode,
        ];
    }

    /**
     * 扫描登录
     * @param $uuid
     * @param string $icon
     * @return array code 408:未扫描;201:扫描未登录;200:登录成功; icon:用户头像
     */
    function checkLogin($uuid, $icon = 'true')
    {
        $url = 'https://login.weixin.qq.com/cgi-bin/mmwebwx-bin/login?loginicon=' . $icon . '&r=' . ~time() . '&uuid=' . $uuid . '&tip=0&_=' . $this->getMillisecond();
        $content = $this->curlPost($url, false, false, 27);
        preg_match('/\d+/', $content, $match);
        $code = $match[0];
        preg_match('/([\'"])([^\'"\.]*?)\1/', $content, $icon);
        if (!empty($icon)) {

            $user_icon = $icon[2];
            $data = array(
                'success' => 1,
                'code' => $code,
                'icon' => $user_icon,
            );
        } else {
            $data['success'] = 1;
            $data['code'] = $code;
        }
        echo json_encode($data);

    }

    /**
     * 登录成功回调
     * @param $uuid
     * @return array $callback
     */
    function get_uri($uuid)
    {
        $url = 'https://login.weixin.qq.com/cgi-bin/mmwebwx-bin/login?uuid=' . $uuid . '&tip=0&_=e' . time();
        $content = $this->curlPost($url);
        $content = explode(';', $content);
        $content_uri = explode('"', $content[1]);
        if (empty($content_uri[1])) {
            exit();
        }
        $uri = $content_uri[1];
        preg_match("~^https:?(//([^/?#]*))?~", $uri, $match);
        $https_header = $match[0];
        $post_url_header = $https_header . "/cgi-bin/mmwebwx-bin";

        $new_uri = explode('scan', $uri);
        $uri = $new_uri[0] . 'fun=new&scan=' . time();
        $getXML = $this->curlPost($uri, false, false, 5, 1);
        $XML = simplexml_load_string($getXML);

        $callback = array(
            'post_url_header' => $post_url_header,
            'Ret' => (array)$XML,
        );
        return $callback;
    }

    /**
     * 获取post数据
     * @param array $callback
     * @return object $this->_response
     */
    function post_self($callback)
    {
        $Ret = $callback['Ret'];
        $status = $Ret['ret'];
        if ($status == '1203') {
            $this->error('未知错误,请2小时后重试');
        }
        if ($status == '0') {
            $this->_response['BaseRequest'] = array(
                'Uin' => $Ret['wxuin'],
                'Sid' => $Ret['wxsid'],
                'Skey' => $Ret['skey'],
                'DeviceID' => 'e' . rand(10000000, 99999999) . rand(1000000, 9999999),
            );

            $this->_response['skey'] = $Ret['skey'];

            $this->_response['pass_ticket'] = $Ret['pass_ticket'];

            $this->_response['sid'] = $Ret['wxsid'];

            $this->_response['uin'] = $Ret['wxuin'];

            $this->_response['header'] = $callback['post_url_header'];

            return $this->_response;
        }
    }

    /**
     * 初始化
     * @return mixed
     */
    function wxinit()
    {
        $url = $this->_response['header'] . '/webwxinit?pass_ticket=' . $this->_response['pass_ticket'] . '&skey=' . $this->_response['skey'] . '&r=' . time();
        $param = array(
            'BaseRequest' => $this->_response['BaseRequest'],
        );
        $json = $this->curlPost($url, $param, false, 35,1);

        $json = json_decode($json, true);

        return $json;
    }

    function getGroupList()
    {
        $data = $this->wxinit();
        $CTL = $data['ContactList'];

        foreach ($CTL as $key => $item) {
            if ($item['MemberCount'] > 0) {
                $GL[] = $item;
            }
        }

        return $GL;
    }

    function getContactList()
    {
        $data = $this->wxinit();
        return $data['ContactList'];
    }

    function getUser()
    {
        $data = $this->wxinit();
        return $data['User'];
    }

    function getSyncKey()
    {
        $data = $this->wxinit();
        return $data['SyncKey'];
    }

    function getUserName($username)
    {
        $this->_response = $_SESSION['post'];

        $url = $this->_response['header'] . '/webwxbatchgetcontact?type=ex&r=' . time() . '&pass_ticket=' . $this->_response[pass_ticket];

        $params = array(
            'BaseRequest' => $this->_response[BaseRequest],
            "Count" => 1,
            "List" => array("UserName" => $username, "EncryChatRoomId" => ""),
        );

        $data = $this->curlPost($url, $params);

        return $data;
    }

    /**
     * 非好友在群组中查找
     * @param $id
     * @param $FromUserName
     * @return string
     */
    function getUserInGroup($id, $FromUserName)
    {

        $GL = $this->getGroupList();

        $user_list = $this->webwxbatchgetcontact($this->_response, $GL);

        foreach ($user_list[ContactList] as $key => $value) {

            if ($value['UserName'] == $FromUserName) {

                $from_list = $value['MemberList'];

                $group_name = $value['NickName'];
            }
        }

        foreach ($from_list as $k => $val) {

            if ($val['UserName'] == $id) {
                $sendmsg_user['DisplayName'] = $val['DisplayName'];
                $sendmsg_user['NickName'] = $val['NickName'];
                $sendmsg_user['GroupName'] = $group_name;
            }
        }

        return $sendmsg_user;
    }

    /**
     * 获取MsgId
     * @return array $data
     */
    function wxstatusnotify()
    {
        $User = $this->getUser();

        $url = $this->_response['header'] . '/webwxstatusnotify?lang=zh_CN&pass_ticket=' . $this->_response['pass_ticket'];

        $params = array(
            'BaseRequest' => $this->_response['BaseRequest'],
            "Code" => 3,
            "FromUserName" => $User['UserName'],
            "ToUserName" => $User['UserName'],
            "ClientMsgId" => time()
        );

        $data = $this->curlPost($url, $params);

        $data = json_decode($data, true);

        return $data;
    }

    /**
     * 获取联系人
     * @return array $data
     */
    function webwxgetcontact()
    {

        $url = $this->_response['header'] . '/webwxgetcontact?pass_ticket=' . $this->_response['pass_ticket'] . '&seq=0&skey=' . $this->_response['skey'] . '&r=' . time();

        $params['BaseRequest'] = $this->_response['BaseRequest'];

        $data = $this->curlPost($url, $params);

        return $data;
    }

    /**
     * 获取当前活跃群信息
     * @param $group_list 从获取联系人和初始化中获取
     * @return array $data
     */
    function webwxbatchgetcontact($group_list)
    {
        $url = $this->_response['header'] . '/webwxbatchgetcontact?type=ex&lang=zh_CN&r=' . time() . '&pass_ticket=' . $this->_response['pass_ticket'];

        $params['BaseRequest'] = $this->_response['BaseRequest'];

        $params['Count'] = count($group_list);

        foreach ($group_list as $key => $value) {
            if ($value[MemberCount] == 0) {
                $params['List'][] = array(
                    'UserName' => $value['UserName'],
                    'ChatRoomId' => "",
                );
            }
            $params['List'][] = array(
                'UserName' => $value['UserName'],
                'EncryChatRoomId' => "",
            );

        }
//        var_dump($params);
        $data = $this->curlPost($url, $params);

        $data = json_decode($data, true);

        return $data;
    }

    public function ImplodeSynckey($synckeylist)
    {
        $str_value = [];
        foreach ($synckeylist as $key => $value) {
            $str_value[] = $value['Key'] . '_' . $value['Val'];
        }
        $sv = urlencode(implode('|', $str_value));
        return $sv;
    }

    /**
     * 心跳检测 0正常；1101失败／登出；2新消息；7不要耍手机了我都收不到消息了；
     * @param $SyncKey 初始化方法中获取
     * @param $timeout
     * @return array $status
     */
    function synccheck($SyncKey)
    {
        if (intval($SyncKey['Count']) > 4) {
            $type = 1;
        } else {
            $type = 0;
        }
        $deviceid = 'e' . rand(10000000, 99999999) . rand(1000000, 9999999);
        $sv = $this->ImplodeSynckey($SyncKey['List']);
        $header = array(
            '0' => 'https://webpush.wx2.qq.com',
            '1' => 'https://webpush.wx.qq.com',
            '2' => 'https://webpush.weixin.qq.com',
            '3' => 'https://webpush2.weixin.qq.com',
            '4' => 'https://webpush.wx8.qq.com',
            '5' => 'https://webpush.web2.wechat.com',
            '6' => 'https://webpush.web.wechat.com',

        );

        foreach ($header as $key => $value) {
            $url = $value . "/cgi-bin/mmwebwx-bin/synccheck?r=" . $this->getMillisecond() . "&skey=" . $this->_response['skey'] . "&sid=" . $this->_response['sid'] . "&deviceid=" . $deviceid . "&uin=" . $this->_response['uin'] . "&synckey=" . $sv . "&_=" . $this->getMillisecond() . "&pass_ticket=" . $this->_response['pass_ticket'];
            $status = $this->curlPost($url, '', false, 30);
            $rule = '/window.synccheck={retcode:"(\d+)",selector:"(\d+)"}/';
            preg_match($rule, $status, $match);
            if ($match) {
                $retcode = intval($match[1]);
                $selector = $match[2];
                if ($retcode == 0) {
                    $arr = array(
                        'ret' => $retcode,
                        'sel' => intval($selector),
                        'key' => $sv,
                    );
                    $model = new Synckey();
                    $model->compare($SyncKey, $sv, $type);
                    return $arr;
                }
            }
        }
        return false;

    }

    /**
     * 获取最新消息
     * @param $SyncKey
     * @return array $data
     */
    function webwxsync($SyncKey)
    {

        $url = $this->_response['header'] . '/webwxsync?sid=' . $this->_response['sid'] . '&skey=' . $this->_response['skey'] . '&pass_ticket=' . $this->_response['pass_ticket'];

        $params = array(
            'BaseRequest' => $this->_response['BaseRequest'],
            'SyncKey' => $SyncKey,
            'rr' => ~time(),
        );

        $data = $this->curlPost($url, $params, $is_gbk = false, $timeout = 35);
        return json_decode($data, true);
    }

    /**
     * 发送消息
     * @param $to 发送人
     * @param $word
     * @return array $data
     */
    function webwxsendmsg($to, $word)
    {

        header("Content-Type: text/html; charset=UTF-8");
        $User = $this->getUser();
        $url = $this->_response['header'] . '/webwxsendmsg?pass_ticket=' . $this->_response[pass_ticket];

        $clientMsgId = getMillisecond() * 1000 + rand(1000, 9999);

        $params = array(
            'BaseRequest' => $this->_response[BaseRequest],
            'Msg' => array(
                "Type" => 1,
                "Content" => $word,
                "FromUserName" => $User['UserName'],
                "ToUserName" => $to,
                "LocalID" => $clientMsgId,
                "ClientMsgId" => $clientMsgId
            ),
            'Scene' => 0,
        );

        $data = $this->curlPost($url, $params, 1);

        return $data;
    }

    function webwxgetimg($MsgID)
    {
        $url = $this->_response['header'] . '/webwxgetmsgimg?&MsgID=' . $MsgID . '&skey=' . urlencode($this->_response['skey']) . '&type=slave';
        return $url;
    }

    function webwxsendimg($to, $MsgId)
    {
        header("Content-Type: text/html; charset=UTF-8");
        $User = $this->getUser();
        $url = $this->_response['header'] . '/webwxsendmsgimg?fun=async&f=json&pass_ticket=' . $this->_response[pass_ticket];

        $clientMsgId = getMillisecond() * 1000 + rand(1000, 9999);

        $params = array(
            'BaseRequest' => $this->_response[BaseRequest],
            'Msg' => array(
                "Type" => 3,
                "MediaId" => $MsgId,
                "FromUserName" => $User['UserName'],
                "ToUserName" => $to,
                "LocalID" => $clientMsgId,
                "ClientMsgId" => $clientMsgId
            ),
            'Scene' => 0,
        );

        $data = $this->curlPost($url, $params, 1, 60);
        return $data;
    }

    /**
     *退出登录
     * @return bool
     */
    function wxloginout()
    {
        $url = $this->_response['header'] . '/webwxlogout?redirect=1&type=1&skey=' . urlencode($this->_response['skey']);
        $param = array(
            'sid' => $this->_response['sid'],
            'uin' => $this->_response['uin'],
        );
        $this->curlPost($url, $param);

        return true;
    }

    /**
     * 抓取https网站爬虫
     * @param $url
     * @param string $data
     * @param bool $is_gbk
     * @param int $timeout
     * @param string $iscookie
     * @param bool $CA
     * @return mixed
     */
    function curlPost($url, $data = '', $is_gbk = false, $timeout = 25, $iscookie = '', $CA = false)
    {
        $cacert = getcwd() . '/cacert.pem'; //CA根证书

        $SSL = substr($url, 0, 8) == "https://" ? true : false;

        $cookie_jar = EXTEND_PATH . '/WxApi/wx.cookie';

        $header[] = "Accept: text/xml,application/xml,application/xhtml+xml,text/html;q=0.9,text/plain;q=0.8,image/png,*/*;q=0.5";
        $header[] = "Cache-Control: max-age=0";
        $header[] = "Connection: keep-alive";
        $header[] = "Keep-Alive: 300";
        $header[] = "Accept-Charset: ISO-8859-1,utf-8;q=0.7,*;q=0.7";
        $header[] = "Accept-Language: en-us,en;q=0.5";
        $header[] = "Pragma: "; // browsers keep this blank.
        $ch = curl_init();

        curl_setopt($ch, CURLOPT_URL, $url);

        curl_setopt($ch, CURLOPT_USERAGENT, 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/54.0.2840.98 Safari/537.36');

        curl_setopt($ch, CURLOPT_HTTPHEADER, $header);

        curl_setopt($ch, CURLOPT_REFERER, 'https://wx.qq.com');

        curl_setopt($ch, CURLOPT_TIMEOUT, $timeout);

        curl_setopt($ch, CURLOPT_CONNECTTIMEOUT, $timeout - 2);
        if ($SSL && $CA) {
            curl_setopt($ch, CURLOPT_SSL_VERIFYPEER, true);   // 只信任CA颁布的证书
            curl_setopt($ch, CURLOPT_CAINFO, $cacert); // CA根证书（用来验证的网站证书是否是CA颁布）
            curl_setopt($ch, CURLOPT_SSL_VERIFYHOST, 2); // 检查证书中是否设置域名，并且是否与提供的主机名匹配
        } else if ($SSL && !$CA) {
            curl_setopt($ch, CURLOPT_SSL_VERIFYPEER, false); // 信任任何证书
            curl_setopt($ch, CURLOPT_SSL_VERIFYHOST, 2); // 检查证书中是否设置域名
        }
        curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
        curl_setopt($ch, CURLOPT_HTTPHEADER, array('Expect:')); //避免data数据过长问题
        if ($data) {
            if ($is_gbk) {
                $data = json_encode($data, JSON_UNESCAPED_UNICODE);
            } else {
                $data = json_encode($data);
            }
            curl_setopt($ch, CURLOPT_POST, true);
            curl_setopt($ch, CURLOPT_POSTFIELDS, $data);
        }
        if ($iscookie > 0) {
            curl_setopt($ch, CURLOPT_COOKIEJAR, $cookie_jar);
        }
        curl_setopt($ch, CURLOPT_COOKIE, EXTEND_PATH . '/WxApi/wx.cookie');
        //curl_setopt($ch, CURLOPT_POSTFIELDS, http_build_query($data)); //data with URLEncode
        $ret = curl_exec($ch);

        curl_close($ch);

        return $ret;
    }

    /**
     * 抓取图片爬虫
     * @param $url
     * @param string $data
     * @param $header
     * @param int $timeout
     * @param bool $CA
     * @return mixed
     */
    function curlFile($url, $data = '', $header, $timeout = 60, $CA = false)
    {
        $cacert = getcwd() . '/cacert.pem'; //CA根证书

        $SSL = substr($url, 0, 8) == "https://" ? true : false;

        $ch = curl_init();
        curl_setopt($ch, CURLOPT_HTTPHEADER, $header);
        curl_setopt($ch, CURLOPT_URL, $url);
        curl_setopt($ch, CURLOPT_REFERER, 'https://wx.qq.com');
        curl_setopt($ch, CURLOPT_TIMEOUT, $timeout);
        curl_setopt($ch, CURLOPT_CONNECTTIMEOUT, $timeout - 2);
        if ($SSL && $CA) {
            curl_setopt($ch, CURLOPT_SSL_VERIFYPEER, true);   // 只信任CA颁布的证书
            curl_setopt($ch, CURLOPT_CAINFO, $cacert); // CA根证书（用来验证的网站证书是否是CA颁布）
            curl_setopt($ch, CURLOPT_SSL_VERIFYHOST, 2); // 检查证书中是否设置域名，并且是否与提供的主机名匹配
        } else if ($SSL && !$CA) {
            curl_setopt($ch, CURLOPT_SSL_VERIFYPEER, false); // 信任任何证书
            curl_setopt($ch, CURLOPT_SSL_VERIFYHOST, true); // 检查证书中是否设置域名
        }
        curl_setopt($ch, CURLOPT_SAFE_UPLOAD, true);
        curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
        curl_setopt($ch, CURLOPT_POSTFIELDS, $data);
        $ret = curl_exec($ch);
        $info = curl_getinfo($ch);

        curl_close($ch);

        return $ret;
    }

    function curlImg($url, $timeout = 60, $CA = false)
    {
        $cacert = getcwd() . '/cacert.pem'; //CA根证书

        $SSL = substr($url, 0, 8) == "https://" ? true : false;
        $header[] = "Accept: text/xml,application/xml,application/xhtml+xml,text/html;q=0.9,text/plain;q=0.8,image/png,*/*;q=0.5";
        $header[] = "Cache-Control: max-age=0";
        $header[] = "Connection: keep-alive";
        $header[] = "Keep-Alive: 300";
        $header[] = "Accept-Charset: ISO-8859-1,utf-8;q=0.7,*;q=0.7";
        $header[] = "Accept-Language: en-us,en;q=0.5";
        $header[] = "Pragma: "; // browsers keep this blank.
        $ch = curl_init();
        $filename = WX_DATA_PATH . date("Ymdhis") . '.jpeg';
        curl_setopt($ch, CURLOPT_HTTPHEADER, $header);
        curl_setopt($ch, CURLOPT_URL, $url);
        curl_setopt($ch, CURLOPT_REFERER, 'https://wx.qq.com');
        curl_setopt($ch, CURLOPT_TIMEOUT, $timeout);
        curl_setopt($ch, CURLOPT_CONNECTTIMEOUT, $timeout - 2);
        if ($SSL && $CA) {
            curl_setopt($ch, CURLOPT_SSL_VERIFYPEER, true);   // 只信任CA颁布的证书
            curl_setopt($ch, CURLOPT_CAINFO, $cacert); // CA根证书（用来验证的网站证书是否是CA颁布）
            curl_setopt($ch, CURLOPT_SSL_VERIFYHOST, 2); // 检查证书中是否设置域名，并且是否与提供的主机名匹配
        } else if ($SSL && !$CA) {
            curl_setopt($ch, CURLOPT_SSL_VERIFYPEER, false); // 信任任何证书
            curl_setopt($ch, CURLOPT_SSL_VERIFYHOST, true); // 检查证书中是否设置域名
        }
        curl_setopt($ch, CURLOPT_SAFE_UPLOAD, true);
        curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
        $ret = curl_exec($ch);
        $info = curl_getinfo($ch);

        curl_close($ch);
        $tp = @fopen($filename, 'a');
        fwrite($tp, $ret);
        fclose($tp);
        return $ret;
    }
}
