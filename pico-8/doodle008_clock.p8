pico-8 cartridge // http://www.pico-8.com
version 16
__lua__
-- wake up clock
-- by: tim swast

--0,1,2,3,4,5,6,7,8,9,10,11
--12,13,14,15

clrramp={0,1,2,8,14,15,7}
xo=0
yo=0
ydir=-1
clkwidth=24
movecounter=0

--darken the whole screen by 1
function darken()
 local x,y,c
 pal()
 for x=0,7 do
  for y=0,7 do
   c=sget(x,y)
   --don't go all to black
   if c>2 then
    c=c-1-flr(rnd(2))
   end
   sset(x,y,c)
  end
 end
end

function dot()
 local x,y,c,glow,tx,ty
 x=flr(rnd(8))
 y=flr(rnd(8))
 pal()
 c=sget(x,y)
 --offset the clock face
 tx=x+xo
 if tx>=clkwidth then
  tx=tx-clkwidth
 elseif tx<0 then
  tx=tx+clkwidth
 end
 ty=y+yo
 if ty>=8 then
  ty=ty-8
 elseif ty<0 then
  ty=ty+8
 end
 glow=sget(tx+8,ty)==1
 if glow then
  c=#clrramp-1
 elseif c>0 then
  c=c-1
 end
 sset(x,y,c)
end

function _update()
 local i
 for i=1,16 do
  dot()
 end
 --move the clock?
 movecounter=movecounter+1
 if movecounter>=30*2 then
  movecounter=0
  darken()
  --move x
  xo=xo+1
  xo=xo%clkwidth
  --move y
  if yo<-2 then
   ydir=1
  elseif yo>-1 then
   ydir=-1
  end
  if rnd(1)<0.1 then
   yo=yo+ydir
  end
 end
end

function _draw()
 local i
 cls()
 spr(0,0,0)
 print(
  stat(93)..":"..stat(94),
  8,0,1)
 --save first row of screen
 --to spritesheet
 memcpy(0x0,0x6000,512)
 --draw big version of screen
 for i=1,#clrramp do
  pal(i-1,clrramp[i])
 end
 palt(0,false)
 sspr(0,0,8,8,0,0,128,128)
end
__label__
00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000
00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000
00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000
00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000
00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000
00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000
00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000
00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000
00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000
00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000
00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000
00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000
00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000
00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000
00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000
00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000
11111111111111111111111111111111222222222222222200000000000000002222222222222222888888888888888800000000000000000000000000000000
11111111111111111111111111111111222222222222222200000000000000002222222222222222888888888888888800000000000000000000000000000000
11111111111111111111111111111111222222222222222200000000000000002222222222222222888888888888888800000000000000000000000000000000
11111111111111111111111111111111222222222222222200000000000000002222222222222222888888888888888800000000000000000000000000000000
11111111111111111111111111111111222222222222222200000000000000002222222222222222888888888888888800000000000000000000000000000000
11111111111111111111111111111111222222222222222200000000000000002222222222222222888888888888888800000000000000000000000000000000
11111111111111111111111111111111222222222222222200000000000000002222222222222222888888888888888800000000000000000000000000000000
11111111111111111111111111111111222222222222222200000000000000002222222222222222888888888888888800000000000000000000000000000000
11111111111111111111111111111111222222222222222200000000000000002222222222222222888888888888888800000000000000000000000000000000
11111111111111111111111111111111222222222222222200000000000000002222222222222222888888888888888800000000000000000000000000000000
11111111111111111111111111111111222222222222222200000000000000002222222222222222888888888888888800000000000000000000000000000000
11111111111111111111111111111111222222222222222200000000000000002222222222222222888888888888888800000000000000000000000000000000
11111111111111111111111111111111222222222222222200000000000000002222222222222222888888888888888800000000000000000000000000000000
11111111111111111111111111111111222222222222222200000000000000002222222222222222888888888888888800000000000000000000000000000000
11111111111111111111111111111111222222222222222200000000000000002222222222222222888888888888888800000000000000000000000000000000
11111111111111111111111111111111222222222222222200000000000000002222222222222222888888888888888800000000000000000000000000000000
77777777777777777777777777777777777777777777777700000000000000007777777777777777777777777777777700000000000000000000000000000000
77777777777777777777777777777777777777777777777700000000000000007777777777777777777777777777777700000000000000000000000000000000
77777777777777777777777777777777777777777777777700000000000000007777777777777777777777777777777700000000000000000000000000000000
77777777777777777777777777777777777777777777777700000000000000007777777777777777777777777777777700000000000000000000000000000000
77777777777777777777777777777777777777777777777700000000000000007777777777777777777777777777777700000000000000000000000000000000
77777777777777777777777777777777777777777777777700000000000000007777777777777777777777777777777700000000000000000000000000000000
77777777777777777777777777777777777777777777777700000000000000007777777777777777777777777777777700000000000000000000000000000000
77777777777777777777777777777777777777777777777700000000000000007777777777777777777777777777777700000000000000000000000000000000
77777777777777777777777777777777777777777777777700000000000000007777777777777777777777777777777700000000000000000000000000000000
77777777777777777777777777777777777777777777777700000000000000007777777777777777777777777777777700000000000000000000000000000000
77777777777777777777777777777777777777777777777700000000000000007777777777777777777777777777777700000000000000000000000000000000
77777777777777777777777777777777777777777777777700000000000000007777777777777777777777777777777700000000000000000000000000000000
77777777777777777777777777777777777777777777777700000000000000007777777777777777777777777777777700000000000000000000000000000000
77777777777777777777777777777777777777777777777700000000000000007777777777777777777777777777777700000000000000000000000000000000
77777777777777777777777777777777777777777777777700000000000000007777777777777777777777777777777700000000000000000000000000000000
77777777777777777777777777777777777777777777777700000000000000007777777777777777777777777777777700000000000000000000000000000000
00000000000000002222222222222222777777777777777700000000000000000000000000000000777777777777777700000000000000000000000000000000
00000000000000002222222222222222777777777777777700000000000000000000000000000000777777777777777700000000000000000000000000000000
00000000000000002222222222222222777777777777777700000000000000000000000000000000777777777777777700000000000000000000000000000000
00000000000000002222222222222222777777777777777700000000000000000000000000000000777777777777777700000000000000000000000000000000
00000000000000002222222222222222777777777777777700000000000000000000000000000000777777777777777700000000000000000000000000000000
00000000000000002222222222222222777777777777777700000000000000000000000000000000777777777777777700000000000000000000000000000000
00000000000000002222222222222222777777777777777700000000000000000000000000000000777777777777777700000000000000000000000000000000
00000000000000002222222222222222777777777777777700000000000000000000000000000000777777777777777700000000000000000000000000000000
00000000000000002222222222222222777777777777777700000000000000000000000000000000777777777777777700000000000000000000000000000000
00000000000000002222222222222222777777777777777700000000000000000000000000000000777777777777777700000000000000000000000000000000
00000000000000002222222222222222777777777777777700000000000000000000000000000000777777777777777700000000000000000000000000000000
00000000000000002222222222222222777777777777777700000000000000000000000000000000777777777777777700000000000000000000000000000000
00000000000000002222222222222222777777777777777700000000000000000000000000000000777777777777777700000000000000000000000000000000
00000000000000002222222222222222777777777777777700000000000000000000000000000000777777777777777700000000000000000000000000000000
00000000000000002222222222222222777777777777777700000000000000000000000000000000777777777777777700000000000000000000000000000000
00000000000000002222222222222222777777777777777700000000000000000000000000000000777777777777777700000000000000000000000000000000
77777777777777777777777777777777777777777777777700000000000000000000000000000000777777777777777700000000000000000000000000000000
77777777777777777777777777777777777777777777777700000000000000000000000000000000777777777777777700000000000000000000000000000000
77777777777777777777777777777777777777777777777700000000000000000000000000000000777777777777777700000000000000000000000000000000
77777777777777777777777777777777777777777777777700000000000000000000000000000000777777777777777700000000000000000000000000000000
77777777777777777777777777777777777777777777777700000000000000000000000000000000777777777777777700000000000000000000000000000000
77777777777777777777777777777777777777777777777700000000000000000000000000000000777777777777777700000000000000000000000000000000
77777777777777777777777777777777777777777777777700000000000000000000000000000000777777777777777700000000000000000000000000000000
77777777777777777777777777777777777777777777777700000000000000000000000000000000777777777777777700000000000000000000000000000000
77777777777777777777777777777777777777777777777700000000000000000000000000000000777777777777777700000000000000000000000000000000
77777777777777777777777777777777777777777777777700000000000000000000000000000000777777777777777700000000000000000000000000000000
77777777777777777777777777777777777777777777777700000000000000000000000000000000777777777777777700000000000000000000000000000000
77777777777777777777777777777777777777777777777700000000000000000000000000000000777777777777777700000000000000000000000000000000
77777777777777777777777777777777777777777777777700000000000000000000000000000000777777777777777700000000000000000000000000000000
77777777777777777777777777777777777777777777777700000000000000000000000000000000777777777777777700000000000000000000000000000000
77777777777777777777777777777777777777777777777700000000000000000000000000000000777777777777777700000000000000000000000000000000
77777777777777777777777777777777777777777777777700000000000000000000000000000000777777777777777700000000000000000000000000000000
eeeeeeeeeeeeeeee1111111111111111111111111111111100000000000000000000000000000000777777777777777711111111111111110000000000000000
eeeeeeeeeeeeeeee1111111111111111111111111111111100000000000000000000000000000000777777777777777711111111111111110000000000000000
eeeeeeeeeeeeeeee1111111111111111111111111111111100000000000000000000000000000000777777777777777711111111111111110000000000000000
eeeeeeeeeeeeeeee1111111111111111111111111111111100000000000000000000000000000000777777777777777711111111111111110000000000000000
eeeeeeeeeeeeeeee1111111111111111111111111111111100000000000000000000000000000000777777777777777711111111111111110000000000000000
eeeeeeeeeeeeeeee1111111111111111111111111111111100000000000000000000000000000000777777777777777711111111111111110000000000000000
eeeeeeeeeeeeeeee1111111111111111111111111111111100000000000000000000000000000000777777777777777711111111111111110000000000000000
eeeeeeeeeeeeeeee1111111111111111111111111111111100000000000000000000000000000000777777777777777711111111111111110000000000000000
eeeeeeeeeeeeeeee1111111111111111111111111111111100000000000000000000000000000000777777777777777711111111111111110000000000000000
eeeeeeeeeeeeeeee1111111111111111111111111111111100000000000000000000000000000000777777777777777711111111111111110000000000000000
eeeeeeeeeeeeeeee1111111111111111111111111111111100000000000000000000000000000000777777777777777711111111111111110000000000000000
eeeeeeeeeeeeeeee1111111111111111111111111111111100000000000000000000000000000000777777777777777711111111111111110000000000000000
eeeeeeeeeeeeeeee1111111111111111111111111111111100000000000000000000000000000000777777777777777711111111111111110000000000000000
eeeeeeeeeeeeeeee1111111111111111111111111111111100000000000000000000000000000000777777777777777711111111111111110000000000000000
eeeeeeeeeeeeeeee1111111111111111111111111111111100000000000000000000000000000000777777777777777711111111111111110000000000000000
eeeeeeeeeeeeeeee1111111111111111111111111111111100000000000000000000000000000000777777777777777711111111111111110000000000000000
77777777777777777777777777777777777777777777777700000000000000007777777777777777777777777777777777777777777777770000000000000000
77777777777777777777777777777777777777777777777700000000000000007777777777777777777777777777777777777777777777770000000000000000
77777777777777777777777777777777777777777777777700000000000000007777777777777777777777777777777777777777777777770000000000000000
77777777777777777777777777777777777777777777777700000000000000007777777777777777777777777777777777777777777777770000000000000000
77777777777777777777777777777777777777777777777700000000000000007777777777777777777777777777777777777777777777770000000000000000
77777777777777777777777777777777777777777777777700000000000000007777777777777777777777777777777777777777777777770000000000000000
77777777777777777777777777777777777777777777777700000000000000007777777777777777777777777777777777777777777777770000000000000000
77777777777777777777777777777777777777777777777700000000000000007777777777777777777777777777777777777777777777770000000000000000
77777777777777777777777777777777777777777777777700000000000000007777777777777777777777777777777777777777777777770000000000000000
77777777777777777777777777777777777777777777777700000000000000007777777777777777777777777777777777777777777777770000000000000000
77777777777777777777777777777777777777777777777700000000000000007777777777777777777777777777777777777777777777770000000000000000
77777777777777777777777777777777777777777777777700000000000000007777777777777777777777777777777777777777777777770000000000000000
77777777777777777777777777777777777777777777777700000000000000007777777777777777777777777777777777777777777777770000000000000000
77777777777777777777777777777777777777777777777700000000000000007777777777777777777777777777777777777777777777770000000000000000
77777777777777777777777777777777777777777777777700000000000000007777777777777777777777777777777777777777777777770000000000000000
77777777777777777777777777777777777777777777777700000000000000007777777777777777777777777777777777777777777777770000000000000000
00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000
00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000
00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000
00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000
00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000
00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000
00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000
00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000
00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000
00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000
00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000
00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000
00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000
00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000
00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000
00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000
