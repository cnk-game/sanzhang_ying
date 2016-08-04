package main

import (
	"log"
	"math/rand"
)

func sortedPlayers() []int {
	offset := 0
	sorted := []int{}
	total := 0

	players := make(map[int]int)
	players[0] = 1
	players[1] = 50
	players[2] = 1

	for i := 0; i < len(players); i++ {
		total = 0
		for k, v := range players {
			exist := false
			for _, pos := range sorted {
				if k == pos {
					exist = true
					break
				}
			}
			if exist {
				continue
			}
			total += v
		}

		//		if total == 0 {
		//			for k := range players {
		//				exist := false
		//				for _, pos := range sorted {
		//					if k == pos {
		//						exist = true
		//						continue
		//					}
		//				}
		//				if !exist {
		//					sorted = append(sorted, k)
		//					break
		//				}
		//			}
		//			log.Println("提前退出sorted:", sorted)
		//			return sorted
		//		}

		r := 0
		if total > 0 {
			r = rand.Int() % total
		}

		log.Println("r = ", r, " total=", total)
		offset = 0
		for k, v := range players {
			exist := false
			for _, pos := range sorted {
				if k == pos {
					exist = true
					continue
				}
			}
			if exist {
				continue
			}

			log.Println("====>r:", r, " offset:", offset, " offset+v:", offset+v)
			if r >= offset && r < offset+v {
				sorted = append(sorted, k)
			}
			offset += v
		}
	}

	return sorted
}

func main() {
	log.Println("排序后:", sortedPlayers())
}
