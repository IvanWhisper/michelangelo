package generate

const (
	/*
		通过数据库记录表Id生成，返回的值是最后的值
		1. 使用事务，保证数据安全，操作隔离
		2. 若先用update，影响行数为0再执行insert，并发下会发生死锁，采用ON DUPLICATE KEY UPDATE实现避免死锁（mysql5.7的坑）
	*/
	SQLInOrUp      = "INSERT INTO idrecord(`Name`,CurValue) VALUE(?,?) ON DUPLICATE KEY UPDATE CurValue = CurValue + ?"
	LuaStInc       = "local temp local source = redis.call('get',KEYS[1]) local num = redis.call('get',KEYS[2]) if source and ((tonumber(num)+ARGV[1])<=(tonumber(source)+ARGV[2])) then temp = redis.call('incrby',KEYS[2],ARGV[1]) return source .. ',' .. tostring(temp) end return '-1,-1'"
	LuaBuildSource = `
local source = redis.call('get',KEYS[1])
if source and (tonumber(source) >= tonumber(ARGV[1]))
then
return 0
else
return redis.call('set',KEYS[1],ARGV[1])
end
`
)
