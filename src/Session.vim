let SessionLoad = 1
let s:so_save = &so | let s:siso_save = &siso | set so=0 siso=0
let v:this_session=expand("<sfile>:p")
silent only
cd ~/go/src/github.com/RenatoGeh/gospn/src
if expand('%') == '' && !&modified && line('$') <= 1 && getline(1) == ''
  let s:wipebuf = bufnr('%')
endif
set shortmess=aoO
badd +142 learn/gens.go
badd +24 utils/indgraph.go
badd +464 main.go
badd +85 io/input.go
badd +110 io/output.go
badd +65 utils/unionfind.go
badd +329 utils/indtest.go
badd +58 utils/kmeans.go
badd +31 spn/sum.go
badd +1 spn/product.go
badd +116 spn/univdist.go
badd +25 spn/node.go
badd +2 spn/varset.go
badd +49 utils/log.go
badd +1 utils/vardata.go
argglobal
silent! argdel *
argadd learn/gens.go
edit utils/indtest.go
set splitbelow splitright
wincmd _ | wincmd |
vsplit
1wincmd h
wincmd _ | wincmd |
split
1wincmd k
wincmd w
wincmd w
wincmd _ | wincmd |
split
1wincmd k
wincmd w
set nosplitbelow
set nosplitright
wincmd t
set winheight=1 winwidth=1
exe '1resize ' . ((&lines * 29 + 30) / 61)
exe 'vert 1resize ' . ((&columns * 104 + 106) / 212)
exe '2resize ' . ((&lines * 28 + 30) / 61)
exe 'vert 2resize ' . ((&columns * 104 + 106) / 212)
exe '3resize ' . ((&lines * 29 + 30) / 61)
exe 'vert 3resize ' . ((&columns * 107 + 106) / 212)
exe '4resize ' . ((&lines * 28 + 30) / 61)
exe 'vert 4resize ' . ((&columns * 107 + 106) / 212)
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
let s:l = 329 - ((0 * winheight(0) + 14) / 29)
if s:l < 1 | let s:l = 1 | endif
exe s:l
normal! zt
329
normal! 05|
lcd ~/go/src/github.com/RenatoGeh/gospn/src
wincmd w
argglobal
edit ~/go/src/github.com/RenatoGeh/gospn/src/utils/indgraph.go
setlocal fdm=manual
setlocal fde=0
setlocal fmr={{{,}}}
setlocal fdi=#
setlocal fdl=0
setlocal fml=1
setlocal fdn=20
setlocal fen
silent! normal! zE
let s:l = 80 - ((23 * winheight(0) + 14) / 28)
if s:l < 1 | let s:l = 1 | endif
exe s:l
normal! zt
80
normal! 0
lcd ~/go/src/github.com/RenatoGeh/gospn/src
wincmd w
argglobal
edit ~/go/src/github.com/RenatoGeh/gospn/src/spn/univdist.go
setlocal fdm=manual
setlocal fde=0
setlocal fmr={{{,}}}
setlocal fdi=#
setlocal fdl=0
setlocal fml=1
setlocal fdn=20
setlocal fen
silent! normal! zE
let s:l = 61 - ((28 * winheight(0) + 14) / 29)
if s:l < 1 | let s:l = 1 | endif
exe s:l
normal! zt
61
normal! 0
lcd ~/go/src/github.com/RenatoGeh/gospn/src
wincmd w
argglobal
edit ~/go/src/github.com/RenatoGeh/gospn/src/utils/log.go
setlocal fdm=manual
setlocal fde=0
setlocal fmr={{{,}}}
setlocal fdi=#
setlocal fdl=0
setlocal fml=1
setlocal fdn=20
setlocal fen
silent! normal! zE
let s:l = 58 - ((27 * winheight(0) + 14) / 28)
if s:l < 1 | let s:l = 1 | endif
exe s:l
normal! zt
58
normal! 033|
lcd ~/go/src/github.com/RenatoGeh/gospn/src
wincmd w
exe '1resize ' . ((&lines * 29 + 30) / 61)
exe 'vert 1resize ' . ((&columns * 104 + 106) / 212)
exe '2resize ' . ((&lines * 28 + 30) / 61)
exe 'vert 2resize ' . ((&columns * 104 + 106) / 212)
exe '3resize ' . ((&lines * 29 + 30) / 61)
exe 'vert 3resize ' . ((&columns * 107 + 106) / 212)
exe '4resize ' . ((&lines * 28 + 30) / 61)
exe 'vert 4resize ' . ((&columns * 107 + 106) / 212)
tabedit ~/go/src/github.com/RenatoGeh/gospn/src/io/output.go
set splitbelow splitright
wincmd _ | wincmd |
vsplit
1wincmd h
wincmd w
set nosplitbelow
set nosplitright
wincmd t
set winheight=1 winwidth=1
exe 'vert 1resize ' . ((&columns * 104 + 106) / 212)
exe 'vert 2resize ' . ((&columns * 107 + 106) / 212)
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
let s:l = 104 - ((38 * winheight(0) + 29) / 58)
if s:l < 1 | let s:l = 1 | endif
exe s:l
normal! zt
104
normal! 03|
lcd ~/go/src/github.com/RenatoGeh/gospn/src
wincmd w
argglobal
edit ~/go/src/github.com/RenatoGeh/gospn/src/io/input.go
setlocal fdm=manual
setlocal fde=0
setlocal fmr={{{,}}}
setlocal fdi=#
setlocal fdl=0
setlocal fml=1
setlocal fdn=20
setlocal fen
silent! normal! zE
let s:l = 153 - ((57 * winheight(0) + 29) / 58)
if s:l < 1 | let s:l = 1 | endif
exe s:l
normal! zt
153
normal! 05|
lcd ~/go/src/github.com/RenatoGeh/gospn/src
wincmd w
exe 'vert 1resize ' . ((&columns * 104 + 106) / 212)
exe 'vert 2resize ' . ((&columns * 107 + 106) / 212)
tabedit ~/go/src/github.com/RenatoGeh/gospn/src/main.go
set splitbelow splitright
wincmd _ | wincmd |
vsplit
1wincmd h
wincmd w
set nosplitbelow
set nosplitright
wincmd t
set winheight=1 winwidth=1
exe 'vert 1resize ' . ((&columns * 104 + 106) / 212)
exe 'vert 2resize ' . ((&columns * 107 + 106) / 212)
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
let s:l = 157 - ((30 * winheight(0) + 29) / 58)
if s:l < 1 | let s:l = 1 | endif
exe s:l
normal! zt
157
normal! 03|
lcd ~/go/src/github.com/RenatoGeh/gospn/src
wincmd w
argglobal
edit ~/go/src/github.com/RenatoGeh/gospn/src/learn/gens.go
setlocal fdm=manual
setlocal fde=0
setlocal fmr={{{,}}}
setlocal fdi=#
setlocal fdl=0
setlocal fml=1
setlocal fdn=20
setlocal fen
silent! normal! zE
let s:l = 198 - ((56 * winheight(0) + 29) / 58)
if s:l < 1 | let s:l = 1 | endif
exe s:l
normal! zt
198
normal! 0
lcd ~/go/src/github.com/RenatoGeh/gospn/src
wincmd w
2wincmd w
exe 'vert 1resize ' . ((&columns * 104 + 106) / 212)
exe 'vert 2resize ' . ((&columns * 107 + 106) / 212)
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
