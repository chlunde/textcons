#!/bin/sh
# https://en.wikipedia.org/wiki/Code_page_437
curl -O https://upload.wikimedia.org/wikipedia/commons/f/f8/Codepage-437.png

# https://doc.pfsense.org/index.php/Installing_pfSense
curl -O https://doc.pfsense.org/images/b/b1/Installer_01_launch_early.png

# http://nerdlypleasures.blogspot.no/2015/04/ibm-character-fonts.html
for path in \
    http://1.bp.blogspot.com/-mt8oWJ49rew/VSHWWsZ4rZI/AAAAAAAACpY/7jSIij-kaV0/s1600/pc0-8x8.png \
    http://1.bp.blogspot.com/-BMG9M6D3bQ4/VSHWWQBXY1I/AAAAAAAACpg/l6WKJyyK2sM/s1600/mda9x14.png \
    http://2.bp.blogspot.com/-mO2qndnI8NI/VSHWTipp1hI/AAAAAAAACoQ/no-CMfAT5Jc/s1600/cga8x8b.png \
    http://3.bp.blogspot.com/-iv76kFtl67w/VSHWTUn5B4I/AAAAAAAACqM/ZipbAfJ7oTc/s1600/cga8x8a.png \
    http://4.bp.blogspot.com/-wEQ2dcm-TuE/VSHWXMMcJdI/AAAAAAAACpk/H3NaqJKGexE/s1600/t1k-8x8-437.png \
    http://1.bp.blogspot.com/-YW67SbP46Zc/VSHWZo20NLI/AAAAAAAACqo/f4NpPnPl21M/s1600/vga8x8.png \
    http://2.bp.blogspot.com/-we3VgQEsB-M/VSHWTg32RAI/AAAAAAAACoM/Opx9wrDi-20/s1600/ega8x14.png \
    http://4.bp.blogspot.com/-UwqqDB0Lx38/VSHWTzpzQRI/AAAAAAAACoY/5SV7psPM_qQ/s1600/ega9x14.png \
    http://3.bp.blogspot.com/-VgN8RlxBx3A/VSHWZbMhcnI/AAAAAAAACqg/TARplBZhEmY/s1600/vga8x16.png \
    http://3.bp.blogspot.com/-bs0_QjenMuE/VSHWaY7CeCI/AAAAAAAACq8/7RxjAsCbIvE/s1600/vga9x16.png \
    http://4.bp.blogspot.com/-K6oItvOIqoU/VSHWXial5TI/AAAAAAAACp0/x5c7PcxErqY/s1600/tga-8x14-437.png \
    http://3.bp.blogspot.com/-WrxCenZe3Q0/VSHWXh_odvI/AAAAAAAACqA/KmQkhp1lqQQ/s1600/tga-8x14-850.png \
    http://1.bp.blogspot.com/-hm77vl68Xao/VSHWUg9V3-I/AAAAAAAACoo/oowbqNysSRc/s1600/isocpi8x16.png
do
    curl -O $path
done
