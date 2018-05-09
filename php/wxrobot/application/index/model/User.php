<?php
/**
 * Created by PhpStorm.
 * User: 聚米粒
 * Date: 2017/5/8
 * Time: 11:02
 */
namespace app\index\model;

use think\Model;

/**
 * @property  nickname
 */
class User extends Model
{
    private $username;
    private $nickname;
    private $icon;

    /**
     * @param $nickname
     * @param $username
     * @param $icon
     * @return false|int|string
     */
    public function postUser($nickname,$username,$icon){
        $data = [
            'nickname'=>$nickname,
            'username'=>$username,
            'icon'=>$icon,
        ];
        $user = $this->where('nickname',$nickname)->find();
        if($user){
            return $this->save($data,['id'=>$user['id']]);
        }else{
            return $this->insert($data);
        }

    }

    public function getUser(){
        $User = $this->order('id desc')->find();
        return $User;
    }
}