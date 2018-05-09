<?php
/**
 * Created by PhpStorm.
 * User: marcelo
 * Date: 2017/5/4
 * Time: 22:24
 */
namespace app\index\model;

use think\Model;

class Synckey extends Model
{
    private $synckey;
    private $str_synckey;

    public function compare($now_synckey,$sv='',$type=0){
        $last_synckey =$this->order('id desc')->find();
        $last_synckey['synckey'] = json_decode($last_synckey['synckey'],true);
        if($last_synckey['synckey']==$now_synckey){
            return $this->update(['str_synckey'=>$sv],['id'=>$last_synckey['id']]);
        }else{
            return $this->insert(['synckey'=>json_encode($now_synckey),'str_synckey'=>$sv,'type'=>$type]);
        }
    }

    public function StrUpdate($str_value,$id){
        return $this->update(['str_synckey'=>$str_value],['id'=>$id]);
    }
    /**
     * @param mixed $synckey
     * @return Synckey
     */
    public function setSynckey($synckey)
    {
        $this->synckey = $synckey;
        return $this;
    }

    /**
     * @param mixed $str_synckey
     * @return Synckey
     */
    public function setStrSynckey($str_synckey)
    {
        $this->str_synckey = $str_synckey;
        return $this;
    }

}