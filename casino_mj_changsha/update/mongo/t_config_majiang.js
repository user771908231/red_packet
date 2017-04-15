/**
 * Created by saplmm on 2017/1/24.
 */

print('=========WECOME==========');
var cursor = db.t_user.find({"nickname": "saplmm"});
printjson(cursor.toArray());
