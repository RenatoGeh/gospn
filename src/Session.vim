let SessionLoad = 1
let s:so_save = &so | let s:siso_save = &siso | set so=0 siso=0
let v:this_session=expand("<sfile>:p")
silent only
cd ~/go/src/github.com/RenatoGeh/gospn/src
if expand('%') == '' && !&modified && line('$') <= 1 && getline(1) == ''
  let s:wipebuf = bufnr('%')
endif
set shortmess=aoO
badd +174 learn/gens.go
badd +87 main.go
badd +27 io/input.go
badd +207 io/output.go
badd +51 utils/unionfind.go
badd +109 spn/sum.go
badd +61 spn/product.go
badd +84 spn/univdist.go
badd +28 spn/node.go
badd +2 spn/varset.go
badd +50 utils/log.go
badd +1 utils/vardata.go
badd +187 utils/indep/indgraph.go
badd +33 utils/indep/fisher.go
badd +29 utils/cluster/dbscan.go
badd +106 utils/cluster/kmeans.go
badd +112 common/queue.go
badd +3 utils/cluster/metrics/euclid.go
badd +18 utils/cluster/optics.go
badd +1 Makefile
badd +21 main_test.go
badd +1 4
badd +1 learn/variable.go
badd +1 out.put
badd +1 common/color.go
badd +411 io/pbm.go
badd +42 io/pgm.go
badd +13 utils/indep/chisq.go
argglobal
silent! argdel *
argadd learn/gens.go
edit spn/univdist.go
set splitbelow splitright
wincmd _ | wincmd |
split
1wincmd k
wincmd _ | wincmd |
vsplit
1wincmd h
wincmd w
wincmd w
wincmd _ | wincmd |
vsplit
1wincmd h
wincmd w
set nosplitbelow
set nosplitright
wincmd t
set winheight=1 winwidth=1
exe '1resize ' . ((&lines * 29 + 30) / 61)
exe 'vert 1resize ' . ((&columns * 119 + 119) / 239)
exe '2resize ' . ((&lines * 29 + 30) / 61)
exe 'vert 2resize ' . ((&columns * 119 + 119) / 239)
exe '3resize ' . ((&lines * 28 + 30) / 61)
exe 'vert 3resize ' . ((&columns * 119 + 119) / 239)
exe '4resize ' . ((&lines * 28 + 30) / 61)
exe 'vert 4resize ' . ((&columns * 119 + 119) / 239)
argglobal
setlocal fdm=manual
setlocal fde=0
setlocal fmr={{{,}}}
setlocal fdi=#
setlocal fdl=0
setlocal fml=1
setlocal fdn=20
setlocal fen
silent! normal! zE
let s:l = 84 - ((28 * winheight(0) + 14) / 29)
if s:l < 1 | let s:l = 1 | endif
exe s:l
normal! zt
84
normal! 0
wincmd w
argglobal
edit spn/node.go
setlocal fdm=manual
setlocal fde=0
setlocal fmr={{{,}}}
setlocal fdi=#
setlocal fdl=0
setlocal fml=1
setlocal fdn=20
setlocal fen
silent! normal! zE
let s:l = 28 - ((27 * winheight(0) + 14) / 29)
if s:l < 1 | let s:l = 1 | endif
exe s:l
normal! zt
28
normal! 04|
wincmd w
argglobal
edit spn/sum.go
setlocal fdm=manual
setlocal fde=0
setlocal fmr={{{,}}}
setlocal fdi=#
setlocal fdl=0
setlocal fml=1
setlocal fdn=20
setlocal fen
silent! normal! zE
let s:l = 109 - ((27 * winheight(0) + 14) / 28)
if s:l < 1 | let s:l = 1 | endif
exe s:l
normal! zt
109
normal! 0
wincmd w
argglobal
edit spn/product.go
setlocal fdm=manual
setlocal fde=0
setlocal fmr={{{,}}}
setlocal fdi=#
setlocal fdl=0
setlocal fml=1
setlocal fdn=20
setlocal fen
silent! normal! zE
let s:l = 61 - ((13 * winheight(0) + 14) / 28)
if s:l < 1 | let s:l = 1 | endif
exe s:l
normal! zt
61
normal! 05|
wincmd w
exe '1resize ' . ((&lines * 29 + 30) / 61)
exe 'vert 1resize ' . ((&columns * 119 + 119) / 239)
exe '2resize ' . ((&lines * 29 + 30) / 61)
exe 'vert 2resize ' . ((&columns * 119 + 119) / 239)
exe '3resize ' . ((&lines * 28 + 30) / 61)
exe 'vert 3resize ' . ((&columns * 119 + 119) / 239)
exe '4resize ' . ((&lines * 28 + 30) / 61)
exe 'vert 4resize ' . ((&columns * 119 + 119) / 239)
tabedit utils/indep/indgraph.go
set splitbelow splitright
wincmd _ | wincmd |
vsplit
1wincmd h
wincmd w
wincmd _ | wincmd |
split
1wincmd k
wincmd w
set nosplitbelow
set nosplitright
wincmd t
set winheight=1 winwidth=1
exe 'vert 1resize ' . ((&columns * 117 + 119) / 239)
exe '2resize ' . ((&lines * 28 + 30) / 61)
exe 'vert 2resize ' . ((&columns * 121 + 119) / 239)
exe '3resize ' . ((&lines * 29 + 30) / 61)
exe 'vert 3resize ' . ((&columns * 121 + 119) / 239)
argglobal
setlocal fdm=manual
setlocal fde=0
setlocal fmr={{{,}}}
setlocal fdi=#
setlocal fdl=0
setlocal fml=1
setlocal fdn=20
setlocal fen
silent! normal! zE
let s:l = 3 - ((2 * winheight(0) + 29) / 58)
if s:l < 1 | let s:l = 1 | endif
exe s:l
normal! zt
3
normal! 0
wincmd w
argglobal
edit utils/indep/chisq.go
setlocal fdm=manual
setlocal fde=0
setlocal fmr={{{,}}}
setlocal fdi=#
setlocal fdl=0
setlocal fml=1
setlocal fdn=20
setlocal fen
silent! normal! zE
let s:l = 360 - ((17 * winheight(0) + 14) / 28)
if s:l < 1 | let s:l = 1 | endif
exe s:l
normal! zt
360
normal! 014|
wincmd w
argglobal
edit utils/indep/fisher.go
setlocal fdm=manual
setlocal fde=0
setlocal fmr={{{,}}}
setlocal fdi=#
setlocal fdl=0
setlocal fml=1
setlocal fdn=20
setlocal fen
silent! normal! zE
let s:l = 3 - ((2 * winheight(0) + 14) / 29)
if s:l < 1 | let s:l = 1 | endif
exe s:l
normal! zt
3
normal! 035|
wincmd w
exe 'vert 1resize ' . ((&columns * 117 + 119) / 239)
exe '2resize ' . ((&lines * 28 + 30) / 61)
exe 'vert 2resize ' . ((&columns * 121 + 119) / 239)
exe '3resize ' . ((&lines * 29 + 30) / 61)
exe 'vert 3resize ' . ((&columns * 121 + 119) / 239)
tabedit io/pgm.go
set splitbelow splitright
wincmd _ | wincmd |
vsplit
1wincmd h
wincmd w
set nosplitbelow
set nosplitright
wincmd t
set winheight=1 winwidth=1
exe 'vert 1resize ' . ((&columns * 117 + 119) / 239)
exe 'vert 2resize ' . ((&columns * 121 + 119) / 239)
argglobal
setlocal fdm=manual
setlocal fde=0
setlocal fmr={{{,}}}
setlocal fdi=#
setlocal fdl=0
setlocal fml=1
setlocal fdn=20
setlocal fen
silent! normal! zE
let s:l = 309 - ((1 * winheight(0) + 29) / 58)
if s:l < 1 | let s:l = 1 | endif
exe s:l
normal! zt
309
normal! 03|
wincmd w
argglobal
edit io/input.go
setlocal fdm=manual
setlocal fde=0
setlocal fmr={{{,}}}
setlocal fdi=#
setlocal fdl=0
setlocal fml=1
setlocal fdn=20
setlocal fen
silent! normal! zE
let s:l = 58 - ((56 * winheight(0) + 29) / 58)
if s:l < 1 | let s:l = 1 | endif
exe s:l
normal! zt
58
normal! 0
wincmd w
exe 'vert 1resize ' . ((&columns * 117 + 119) / 239)
exe 'vert 2resize ' . ((&columns * 121 + 119) / 239)
tabedit main.go
set splitbelow splitright
wincmd _ | wincmd |
vsplit
1wincmd h
wincmd w
wincmd _ | wincmd |
split
1wincmd k
wincmd w
set nosplitbelow
set nosplitright
wincmd t
set winheight=1 winwidth=1
exe 'vert 1resize ' . ((&columns * 117 + 119) / 239)
exe '2resize ' . ((&lines * 29 + 30) / 61)
exe 'vert 2resize ' . ((&columns * 121 + 119) / 239)
exe '3resize ' . ((&lines * 28 + 30) / 61)
exe 'vert 3resize ' . ((&columns * 121 + 119) / 239)
argglobal
setlocal fdm=manual
setlocal fde=0
setlocal fmr={{{,}}}
setlocal fdi=#
setlocal fdl=0
setlocal fml=1
setlocal fdn=20
setlocal fen
silent! normal! zE
let s:l = 15 - ((14 * winheight(0) + 29) / 58)
if s:l < 1 | let s:l = 1 | endif
exe s:l
normal! zt
15
normal! 027|
wincmd w
argglobal
edit main.go
setlocal fdm=manual
setlocal fde=0
setlocal fmr={{{,}}}
setlocal fdi=#
setlocal fdl=0
setlocal fml=1
setlocal fdn=20
setlocal fen
silent! normal! zE
let s:l = 87 - ((14 * winheight(0) + 14) / 29)
if s:l < 1 | let s:l = 1 | endif
exe s:l
normal! zt
87
normal! 05|
wincmd w
argglobal
edit io/pgm.go
setlocal fdm=manual
setlocal fde=0
setlocal fmr={{{,}}}
setlocal fdi=#
setlocal fdl=0
setlocal fml=1
setlocal fdn=20
setlocal fen
silent! normal! zE
let s:l = 42 - ((20 * winheight(0) + 14) / 28)
if s:l < 1 | let s:l = 1 | endif
exe s:l
normal! zt
42
normal! 0
wincmd w
exe 'vert 1resize ' . ((&columns * 117 + 119) / 239)
exe '2resize ' . ((&lines * 29 + 30) / 61)
exe 'vert 2resize ' . ((&columns * 121 + 119) / 239)
exe '3resize ' . ((&lines * 28 + 30) / 61)
exe 'vert 3resize ' . ((&columns * 121 + 119) / 239)
tabedit common/color.go
set splitbelow splitright
set nosplitbelow
set nosplitright
wincmd t
set winheight=1 winwidth=1
argglobal
setlocal fdm=manual
setlocal fde=0
setlocal fmr={{{,}}}
setlocal fdi=#
setlocal fdl=0
setlocal fml=1
setlocal fdn=20
setlocal fen
silent! normal! zE
let s:l = 2 - ((1 * winheight(0) + 29) / 58)
if s:l < 1 | let s:l = 1 | endif
exe s:l
normal! zt
2
normal! 0
tabedit learn/gens.go
set splitbelow splitright
wincmd _ | wincmd |
vsplit
1wincmd h
wincmd w
wincmd _ | wincmd |
split
1wincmd k
wincmd w
set nosplitbelow
set nosplitright
wincmd t
set winheight=1 winwidth=1
exe 'vert 1resize ' . ((&columns * 119 + 119) / 239)
exe '2resize ' . ((&lines * 29 + 30) / 61)
exe 'vert 2resize ' . ((&columns * 119 + 119) / 239)
exe '3resize ' . ((&lines * 28 + 30) / 61)
exe 'vert 3resize ' . ((&columns * 119 + 119) / 239)
argglobal
setlocal fdm=manual
setlocal fde=0
setlocal fmr={{{,}}}
setlocal fdi=#
setlocal fdl=0
setlocal fml=1
setlocal fdn=20
setlocal fen
silent! normal! zE
let s:l = 174 - ((28 * winheight(0) + 29) / 58)
if s:l < 1 | let s:l = 1 | endif
exe s:l
normal! zt
174
normal! 036|
wincmd w
argglobal
edit utils/cluster/metrics/euclid.go
setlocal fdm=manual
setlocal fde=0
setlocal fmr={{{,}}}
setlocal fdi=#
setlocal fdl=0
setlocal fml=1
setlocal fdn=20
setlocal fen
silent! normal! zE
let s:l = 13 - ((12 * winheight(0) + 14) / 29)
if s:l < 1 | let s:l = 1 | endif
exe s:l
normal! zt
13
normal! 0
wincmd w
argglobal
edit utils/cluster/dbscan.go
setlocal fdm=manual
setlocal fde=0
setlocal fmr={{{,}}}
setlocal fdi=#
setlocal fdl=0
setlocal fml=1
setlocal fdn=20
setlocal fen
silent! normal! zE
let s:l = 2 - ((1 * winheight(0) + 14) / 28)
if s:l < 1 | let s:l = 1 | endif
exe s:l
normal! zt
2
normal! 0
wincmd w
exe 'vert 1resize ' . ((&columns * 119 + 119) / 239)
exe '2resize ' . ((&lines * 29 + 30) / 61)
exe 'vert 2resize ' . ((&columns * 119 + 119) / 239)
exe '3resize ' . ((&lines * 28 + 30) / 61)
exe 'vert 3resize ' . ((&columns * 119 + 119) / 239)
tabedit Makefile
set splitbelow splitright
set nosplitbelow
set nosplitright
wincmd t
set winheight=1 winwidth=1
argglobal
setlocal fdm=manual
setlocal fde=0
setlocal fmr={{{,}}}
setlocal fdi=#
setlocal fdl=0
setlocal fml=1
setlocal fdn=20
setlocal fen
silent! normal! zE
let s:l = 21 - ((20 * winheight(0) + 29) / 58)
if s:l < 1 | let s:l = 1 | endif
exe s:l
normal! zt
21
normal! 028|
tabnext 3
if exists('s:wipebuf') && getbufvar(s:wipebuf, '&buftype') isnot# 'terminal'
  silent exe 'bwipe ' . s:wipebuf
endif
unlet! s:wipebuf
set winheight=1 winwidth=20 shortmess=filnxtToO
let s:sx = expand("<sfile>:p:r")."x.vim"
if file_readable(s:sx)
  exe "source " . fnameescape(s:sx)
endif
let &so = s:so_save | let &siso = s:siso_save
let g:this_session = v:this_session
let g:this_obsession = v:this_session
let g:this_obsession_status = 2
doautoall SessionLoadPost
unlet SessionLoad
" vim: set ft=vim :
