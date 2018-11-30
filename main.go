package main

import (
	"fmt"
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

//取结点归属表
func (c *Consistent) SearchTable(node Node) uint32 {
	for i, ring := range c.HashRing {
		if ring >= node.Ring {
			return ring
		}
		if i == c.HashRing.Len()-1 { //如果最后一个table都比hash(node)小，则取第一个table
			return c.HashRing[0]
		}
	}
	return c.HashRing[0]
}

//扩容（表增加）,并迁移相关数据
func (c *Consistent) AddTable(table Table) {
	c.TableMap[table.Ring] = table
	c.HashRing = append(c.HashRing, table.Ring)
	sort.Sort(c.HashRing)
	//找到新增的表的hash在环上的位置,以及新增表hash下一个表hash位置
	var preRing uint32
	for i, ring := range c.HashRing {
		if ring == table.Ring && i != len(c.HashRing)-1 {
			preRing = c.HashRing[i+1]
			fmt.Println("新加入的表所处位置：" + strconv.Itoa(i+1))
			break
		} else if ring == table.Ring && i == len(c.HashRing)-1 { //环状，最后一个hash值下一个就是第一个hash
			fmt.Println("新加入的表所处位置：last")
			preRing = c.HashRing[0]
		}
	}
	preTable := c.TableMap[preRing]
	fmt.Println("迁移数据所在老表：" + preTable.Name)
	fmt.Println("迁移数据所在新表：" + table.Name)
	for _, node := range c.NodeList {
		//迁移部分数据到新表
		//当新的表处于环的第一个table结点
		if table.Ring == c.HashRing[0] {
			//取环上 大于最后一个表hash，或者小于第一个表hash区间的数据
			if node.TableName == preTable.Name && (node.Ring <= table.Ring || node.Ring > c.HashRing[len(c.HashRing)-1]) {
				node.TableName = table.Name
				fmt.Println("迁移数据：" + node.Data)
			}
			continue
		}
		if node.TableName == preTable.Name && node.Ring <= table.Ring {
			node.TableName = table.Name
			fmt.Println("迁移数据：" + node.Data)
		}
	}
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
		selectRing := c.SearchTable(*node)
		node.TableName = c.TableMap[selectRing].Name
		c.NodeList = append(c.NodeList, *node)
	}

	//添加表
	addTable := &Table{}
	addTable.Name = "table-add4"
	addTable.Ring = HashKey(addTable.Name)
	c.AddTable(*addTable) //for _, n := range c.NodeList {
	//	fmt.Println(n.TableName)
	//}

}
