package userService

import "sync"


var UserLockPools UserLockPool

func init(){
	userLockPool := &UserLockPool{}
	userLockPool.pool = make(map[uint32]*UserLock)
}

/**
	用户锁
 */
type UserLock struct {
	sync.Mutex
	userId uint32
}

/**
	存放userLock 相关的东西
 */

type UserLockPool struct {
	pool map[uint32]*UserLock
}


/**
	通过UserId活的用户锁
 */
func (u *UserLockPool) GetUserLockByUserId(userId uint32) *UserLock{
	result := u.pool[userId]
	return result
}


/**
	为用户锁池中增加锁
 */
func (u *UserLockPool) addUserLockByUserId(userId uint32) (*UserLock,error){
	//首先判断pool中是否已经存在,如果存在返回保存失败,如果不存在则从新生成并且返回结果
	result := &UserLock{}
	result.userId = userId
	u.pool[userId] = result
	return result,nil
}

/**
	删除锁池子中的用户锁
 */
func (u *UserLockPool) rmUserLockByUserId(userId uint32) error{
	delete(u.pool,userId);
	return nil
}