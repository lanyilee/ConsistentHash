package main

import (
	"hash/crc32"
	"sort"
	"strconv"
)

//一致性hash算法
//初定10个分表

type Ring []uint32

type Table struct {
	Ring uint32 //此表对应在hash环上的位置
	Name string //表名
}

//结点
type Node struct {
	Ring      uint32 //其hash决定在哪个table
	Data      string //字符串
	TableName string //当前结点被分配的表，在hash环顺时针最近一个table
}

type Consistent struct {
	TableMap map[uint32]Table //所有表详细信息
	HashRing Ring             //所有表的hash
	NodeList []Node           //所有node点信息，在具体项目可以存在数据库中
}

func (c Ring) Len() int {
	return len(c)
}
func (c Ring) Less(i int, j int) bool {
	return c[i] < c[j]
}
func (c Ring) Swap(i int, j int) {
	c[i], c[j] = c[j], c[i]
}

func HashKey(key string) uint32 {
	return crc32.ChecksumIEEE([]byte(key))
}
func (c *Consistent) SearchTable(node Node) uint32 {
	for ring := range c.HashRing {
		if ring > int(node.Ring) {

		}
	}
	return 1
}

func main() {
	c := &Consistent{}
	c.TableMap = make(map[uint32]Table)
	//var HashRing []uint32
	for i := 1; i <= 10; i++ {
		table := &Table{}
		table.Name = "table" + strconv.Itoa(i)
		table.Ring = HashKey(table.Name)
		c.HashRing = append(c.HashRing, table.Ring)
		c.TableMap[table.Ring] = *table
	}
	//将ring排序
	sort.Sort(c.HashRing)
	//分配100个数据到10张表
	for i := 0; i < 100; i++ {
		node := &Node{}
		node.Data = "data" + strconv.Itoa(i)
		node.Ring = HashKey(node.Data)
		//找hash环中距离node最近的table

		c.NodeList = append(c.NodeList, *node)
	}

}
