package ks0108

func abs(x int) int {
	if x < 0 {
		return -x;
	}
	return x;
}

// DrawRect - Draws a rect
func (lcd *Ks0108) DrawRect(x int, y int, w int, h int, fill bool){
	if fill {
		for nx:=x; nx < (x+w) ; nx++ {
			for ny:=y; ny < (y+h) ; ny++ {
				lcd.setPixel(uint8(nx),uint8(ny));
			}
		}
	} else {
		lcd.DrawLine(x,y,(x+w)-1,y);
		lcd.DrawLine(x,(y+h)-1,(x+w)-1,(y+h)-1);
		lcd.DrawLine(x,y,x,(y+h)-1);
		lcd.DrawLine((x+w)-1,y,(x+w)-1,(y+h)-1);
	}
}

// DrawLine - Draws a 1 pixel width line from x to y
func (lcd * Ks0108) DrawLine(x0 int, y0 int, x1 int, y1 int){
	var e2 int
	dx := abs(x1-x0);
	dy := abs(y1-y0)
	err := map[bool]int{true: int(dx), false: int(-dx)} [dx > dy];
	sx := map[bool]int{true: 1, false: -1} [x0 < x1];
	sy := map[bool]int{true: 1, false: -1} [y0 < y1];

	for ;; {
		lcd.setPixel(uint8(x0),uint8(y0));
		if x0==x1 && y0==y1 {
			break;
		}
		e2 = err;
		if (e2 >-dx) { err -= dy; x0 += sx; }
		if (e2 < dy) { err += dx; y0 += sy; }
	}

}