package main

import (
	"fmt"
	"math"
)

// method 1: Triangulation
func polygonArea1(x, y []float64) float64 {
	n := len(x)
	area := 0.0
	for i := 0; i < n; i++ {
		j := (i + 1) % n
		area += x[i] * y[j] - x[j] * y[i]
	}
	area = math.Abs(area) / 2.0

	// Triangulate the polygon
	triArea := 0.0
	for i := 1; i < n-1; i++ {
		triArea += triangleArea1(x[0], y[0], x[i], y[i], x[i+1], y[i+1])
	}

	return triArea
}

func triangleArea1(x1, y1, x2, y2, x3, y3 float64) float64 {
	a := math.Sqrt((x1-x2)*(x1-x2) + (y1-y2)*(y1-y2))
	b := math.Sqrt((x2-x3)*(x2-x3) + (y2-y3)*(y2-y3))
	c := math.Sqrt((x3-x1)*(x3-x1) + (y3-y1)*(y3-y1))
	s := (a + b + c) / 2.0
	return math.Sqrt(s * (s - a) * (s - b) * (s - c))
}

// method 2: Gaussian
func polygonArea2(x, y []float64) float64 {
	n := len(x)
	area := 0.0
	for i := 0; i < n; i++ {
		j := (i + 1) % n
		area += x[i] * y[j] - x[j] * y[i]
	}
	return 0.5 * area
}

// method3: Scanline
type point struct {
	x float64
	y float64
}

func polygonArea3(vertices []point) float64 {
	n := len(vertices)
	if n < 3 {
		return 0
	}

	// Sort vertices by y-coordinate
	for i := 0; i < n-1; i++ {
		for j := i + 1; j < n; j++ {
			if vertices[i].y > vertices[j].y {
				vertices[i], vertices[j] = vertices[j], vertices[i]
			}
		}
	}

	// Compute area using scanline algorithm
	area := 0.0
	for y := vertices[0].y + 1; y <= vertices[n-1].y; y++ {
		var intersections []float64
		for i := 0; i < n; i++ {
			j := (i + 1) % n
			if (vertices[i].y <= y && y < vertices[j].y) || (vertices[j].y <= y && y < vertices[i].y) {
				x := (y-vertices[i].y)*(vertices[j].x-vertices[i].x)/(vertices[j].y-vertices[i].y) + vertices[i].x
				intersections = append(intersections, x)
			}
		}
		intersections = append(intersections, math.Inf(1))

		// Sort intersections by x-coordinate
		for i := 0; i < len(intersections)-1; i++ {
			for j := i + 1; j < len(intersections); j++ {
				if intersections[i] > intersections[j] {
					intersections[i], intersections[j] = intersections[j], intersections[i]
				}
			}
		}

		// Compute area contribution of current scanline
		for i := 0; i < len(intersections)-1; i++ {
			x1 := intersections[i]
			x2 := intersections[i+1]
			area += 0.5 * (y - vertices[0].y) * (x2 - x1)
		}
	}

	return area
}

// method4: Divide and Conquer

func PolygonArea4(vertices []point) float64 {
	if len(vertices) < 3 {
		return 0.0
	}

	return dividePolygonArea(vertices, 0, len(vertices)-1)
}

func dividePolygonArea(vertices []point, start int, end int) float64 {
	if end-start == 2 {
		// 三角形
		return triangleArea4(vertices[start], vertices[start+1], vertices[end])
	} else if end-start == 3 {
		// 梯形和三角形
		return trapezoidArea4(vertices[start], vertices[start+1], vertices[start+2], vertices[end])
	}

	// 将多边形分割成两部分
	mid := (start + end) / 2
	leftArea := dividePolygonArea(vertices, start, mid)
	rightArea := dividePolygonArea(vertices, mid, end)

	return leftArea + rightArea
}

func triangleArea4(p1 point, p2 point, p3 point) float64 {
	return 0.5 * (p1.x*p2.y + p2.x*p3.y + p3.x*p1.y - p1.y*p2.x - p2.y*p3.x - p3.y*p1.x)
}

func trapezoidArea4(p1 point, p2 point, p3 point, p4 point) float64 {
	return 0.5 * (p2.y-p1.y + p3.y-p4.y) * (p4.x-p1.x)
}


//method5:  Cut-and-stitch
func PolygonArea5(vertices []point) float64 {
	if len(vertices) < 3 {
		return 0.0
	}

	// 计算多边形重心
	var centroidx, centroidy float64
	for i := 0; i < len(vertices); i++ {
		next := (i + 1) % len(vertices)
		cross := vertices[i].x*vertices[next].y - vertices[next].x*vertices[i].y
		centroidx += (vertices[i].x + vertices[next].x) * cross
		centroidy += (vertices[i].y + vertices[next].y) * cross
	}
	area := PolygonSignedArea5(vertices)
	centroidx /= 6 * area
	centroidy /= 6 * area

	// 计算每个三角形的面积
	totalArea := 0.0
	for i := 0; i < len(vertices); i++ {
		next := (i + 1) % len(vertices)
		area := 0.5 * math.Abs((vertices[i].x-centroidx)*(vertices[next].y-centroidy) - (vertices[next].x-centroidx)*(vertices[i].y-centroidy))
		totalArea += area
	}

	return totalArea
}

func PolygonSignedArea5(vertices []point) float64 {
	if len(vertices) < 3 {
		return 0.0
	}

	sum := 0.0
	for i := 0; i < len(vertices); i++ {
		next := (i + 1) % len(vertices)
		sum += vertices[i].x * vertices[next].y - vertices[next].x * vertices[i].y
	}

	return 0.5 * sum
}


// method6: Vertex Reordering
func VertexReordering(vertices []point) []point {
	// 按照经纬度排序
	vertices = sortByLongitude(vertices)

	// 确定原初标准的步长
	step := computeStep(vertices)

	// 取出所有顶点的最高点和最低点，并进行编号
	top, bottom := getTopAndBottompoints(vertices)
	_, bottomIndex := getTopAndBottomIndices(vertices, top, bottom)

	// 推算出每个顶点所在的地方，并将其分别编号
	columnCount := int(math.Ceil((top.y - bottom.y) / step))
	vertexColumns := make([][]int, columnCount)
	for i := 0; i < len(vertices); i++ {
		column := int(math.Floor((vertices[i].y - bottom.y) / step))
		vertexColumns[column] = append(vertexColumns[column], i)
	}

	// 按照编号的顺序进行推步聚顶计算
	triangles := make([][3]int, 0)
	for i := 0; i < columnCount-1; i++ {
		column1 := vertexColumns[i]
		column2 := vertexColumns[i+1]
		for j := 0; j < len(column1); j++ {
			index1 := column1[j]
			for k := 0; k < len(column2); k++ {
				index2 := column2[k]
				if isDiagonal(vertices, index1, index2) {
					triangles = append(triangles, [3]int{index1, index2, bottomIndex + i})
					triangles = append(triangles, [3]int{index2, index2, bottomIndex + i})
					triangles = append(triangles, [3]int{index2, index1, bottomIndex + i + 1})
					break
				}
			}
		}
	}

	return vertices
}

func sortByLongitude(vertices []point) []point {
	// TODO: 根据经度排序
	return vertices
}

func computeStep(vertices []point) float64 {
	// TODO: 确定原初标准的步长
	return 0.1
}

func getTopAndBottompoints(vertices []point) (point, point) {
	// TODO: 取出最高点和最低点
	return point{}, point{}
}

func getTopAndBottomIndices(vertices []point, top point, bottom point) (int, int) {
	// TODO: 取出最高点和最低点的编号
	return 0, 0
}

func isDiagonal(vertices []point, i int, j int) bool {
	// TODO: 判断是否为多边形的对角线
	return false
}

// main
func main() {
	x := []float64{0, 1, 2, 1}
	y := []float64{0, 0, 1, 2}
	area := polygonArea1(x, y)

	//vertices := []point{{0, 0}, {1, 0}, {2, 1}, {1, 2}}
	//area := PolygonArea(vertices)
	fmt.Printf("The area of the polygon is %f\n", area)
}
