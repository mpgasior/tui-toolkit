package view

import "sort"

type direction int

const (
	dirHorizontal direction = iota
	dirVertical
)

type Constraint struct {
	Name     string
	IsStatic bool
	Value    int
}

func Fixed(name string, value int) Constraint {
	return Constraint{
		Name: name, IsStatic: true, Value: value,
	}
}

func Dynamic(name string, value int) Constraint {
	return Constraint{
		Name: name, IsStatic: false, Value: value,
	}
}

func Grid(p Port, rowConstraints []Constraint, colConstraints []Constraint) map[string]Port {
	result := make(map[string]Port)

	rows := SplitH(p, rowConstraints...)

	for _, rc := range rowConstraints {
		if rc.Name == "" {
			continue
		}

		rowPort := rows[rc.Name]
		cols := SplitV(rowPort, colConstraints...)

		for _, cc := range colConstraints {
			if cc.Name == "" {
				continue
			}

			key := rc.Name + ":" + cc.Name
			result[key] = cols[cc.Name]
		}
	}

	return result
}

func Center(p Port, cw Constraint, ch Constraint) Port {
	centerV := CenterV(p, cw)

	return CenterH(centerV, ch)
}

func CenterH(p Port, cs Constraint) Port {
	split := SplitH(p, Dynamic("", 1), cs, Dynamic("", 1))
	return split[cs.Name]
}

func CenterV(p Port, cs Constraint) Port {
	split := SplitV(p, Dynamic("", 1), cs, Dynamic("", 1))
	return split[cs.Name]
}

func SplitH(p Port, cs ...Constraint) map[string]Port {
	return splitMap(p, dirHorizontal, cs...)
}

func SplitV(p Port, cs ...Constraint) map[string]Port {
	return splitMap(p, dirVertical, cs...)
}

func splitMap(p Port, dir direction, cs ...Constraint) map[string]Port {
	fixedTotal := 0
	weights := make([]int, len(cs))
	for i, c := range cs {
		if c.IsStatic {
			fixedTotal += c.Value
		} else {
			weights[i] = c.Value
		}
	}

	total := p.w
	if dir == dirHorizontal {
		total = p.h
	}

	allocations := distribute(total-fixedTotal, weights...)
	result := make(map[string]Port)

	offset := 0
	for i, allocation := range allocations {
		c := cs[i]

		if c.IsStatic {
			allocation = c.Value
		}

		x, y, w, h := offset, 0, allocation, p.h
		if dir == dirHorizontal {
			x, y, w, h = 0, offset, p.w, allocation
		}

		result[c.Name] = p.Slice(x, y, w, h)
		offset += allocation
	}

	return result
}

func distribute(total int, weights ...int) []int {
	type allocation struct {
		index    int
		value    int
		fraction int
	}

	sumWeights := 0
	for _, w := range weights {
		sumWeights += w
	}

	if sumWeights == 0 {
		return make([]int, len(weights))
	}

	allocations := make([]allocation, len(weights))
	distributed := 0

	for i, w := range weights {
		value := (total * w) / sumWeights
		fraction := (total * w) % sumWeights

		allocations[i] = allocation{
			index:    i,
			value:    value,
			fraction: fraction,
		}

		distributed += value
	}

	leftover := total - distributed

	sort.Slice(allocations, func(i, j int) bool {
		return allocations[i].fraction > allocations[j].fraction
	})

	for i := range leftover {
		allocations[i].value += 1
	}

	result := make([]int, len(weights))
	for _, allocation := range allocations {
		result[allocation.index] = allocation.value
	}
	return result
}
