package main

import (
	"fmt"
	"hash/crc32"
	"sort"
	"strconv"
)

//consistent hash algorithm
// initial 10 tables

type Ring []uint32

type Table struct {
	Ring uint32 // the corresponding position of the table on the ring
	Name string
}

//结点
type Node struct {
	Ring      uint32 //hash(data)
	Data      string //data
	TableName string //the table to which the node belongs to
}

type Consistent struct {
	TableMap map[uint32]Table //all tables message
	HashRing Ring             //all tables hash
	NodeList []Node           //all node message，should be stored on the database if possible
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

//find the table which the given node  belongs to
func (c *Consistent) SearchTable(node Node) uint32 {
	for i, ring := range c.HashRing {
		if ring >= node.Ring {
			return ring
		}
		if i == c.HashRing.Len()-1 { //if the last one table (the largest uint32 num) is less than hash(node),then select the first table(the smallest uint32 num) on the ring
			return c.HashRing[0]
		}
	}
	return c.HashRing[0]
}

//capacity(add table),and migrating data
func (c *Consistent) AddTable(table Table) {
	c.TableMap[table.Ring] = table
	c.HashRing = append(c.HashRing, table.Ring)
	sort.Sort(c.HashRing)
	//find the new table position on the ring
	var preRing uint32
	for i, ring := range c.HashRing {
		if ring == table.Ring && i != len(c.HashRing)-1 {
			preRing = c.HashRing[i+1]
			fmt.Println("added table old position：" + strconv.Itoa(i+1))
			break
		} else if ring == table.Ring && i == len(c.HashRing)-1 { //ring structure，the largest hash(table) next is the smallest hash(table)
			fmt.Println("added table old position：last")
			preRing = c.HashRing[0]
		}
	}
	preTable := c.TableMap[preRing]
	fmt.Println("migrating data is in the old table ：" + preTable.Name)
	fmt.Println("migrating data is in the new table：" + table.Name)
	for _, node := range c.NodeList {
		//the position where  the new table is located on is the smallest hash(new table)
		if table.Ring == c.HashRing[0] {
			if node.TableName == preTable.Name && (node.Ring <= table.Ring || node.Ring > c.HashRing[len(c.HashRing)-1]) {
				node.TableName = table.Name
				fmt.Println("migrating data：" + node.Data)
			}
			continue
		}
		if node.TableName == preTable.Name && node.Ring <= table.Ring {
			node.TableName = table.Name
			fmt.Println("migrating data：" + node.Data)
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
	//sort the ring,from small to large
	sort.Sort(c.HashRing)
	//allocate 100 data to 10 tables
	for i := 0; i < 100; i++ {
		node := &Node{}
		node.Data = "data" + strconv.Itoa(i)
		node.Ring = HashKey(node.Data)
		selectRing := c.SearchTable(*node)
		node.TableName = c.TableMap[selectRing].Name
		c.NodeList = append(c.NodeList, *node)
	}

	//add table
	addTable := &Table{}
	addTable.Name = "table-add4"
	addTable.Ring = HashKey(addTable.Name)
	c.AddTable(*addTable) //for _, n := range c.NodeList {
	//	fmt.Println(n.TableName)
	//}

}
