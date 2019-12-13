package ks0108

// void Ks0108pi::drawRect(uint8_t x, uint8_t y, uint8_t w, uint8_t h, uint8_t style){
// 	for(int nx=x; nx < (x+w) ; nx++){
// 		for(int ny=y; ny < (y+h) ; ny++){
// 			if(style & STYLE_BLACK_BG) setPixel(nx,ny);
// 			else if(style & STYLE_WHITE_BG) clearPixel(nx,ny);
// 			else if(style & STYLE_GRAY_BG){
// 				if((nx+ny)%2==0)
// 					clearPixel(nx,ny);
// 				else
// 					setPixel(nx,ny);
// 			}
// 		}
// 	}

// 	if( (style & STYLE_BLACK_BORDER) || (style & STYLE_WHITE_BORDER)) {
// 		drawLine(x,y,(x+w)-1,y);
// 		drawLine(x,(y+h)-1,(x+w)-1,(y+h)-1);
// 		drawLine(x,y,x,(y+h)-1);
// 		drawLine((x+w)-1,y,(x+w)-1,(y+h)-1);
// 	}
// }

// void Ks0108pi::drawLine(uint8_t x0, uint8_t y0, uint8_t x1, uint8_t y1){
// 	int dx = abs(x1-x0), sx = x0<x1 ? 1 : -1;
// 	int dy = abs(y1-y0), sy = y0<y1 ? 1 : -1;
// 	int err = (dx>dy ? dx : -dy)/2, e2;

// 	for(;;){
// 		setPixel(x0,y0);
// 		if (x0==x1 && y0==y1) break;
// 		e2 = err;
// 		if (e2 >-dx) { err -= dy; x0 += sx; }
// 		if (e2 < dy) { err += dx; y0 += sy; }
// 	}

// }