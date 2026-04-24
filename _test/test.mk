.PHONY: test-canvas
test-canvas: compile # create a canvas with the given size and colour
	@OVERLAY_LOG_LEVEL=d dist/overlay_linux_amd64_v1/overlay canvas --size=640,480 --colour=#FF0000 --output=dist/overlay_linux_amd64_v1/red.png
	@OVERLAY_LOG_LEVEL=d dist/overlay_linux_amd64_v1/overlay canvas --size=640,480 --colour=#00FF00 --output=dist/overlay_linux_amd64_v1/green.jpg
	@OVERLAY_LOG_LEVEL=d dist/overlay_linux_amd64_v1/overlay canvas --size=640,480 --colour=#0000FF --output=dist/overlay_linux_amd64_v1/blue.bmp

.PHONY: test-rectangle
test-rectangle: compile # create a rectangle with the given size and colour
	@OVERLAY_LOG_LEVEL=d dist/overlay_linux_amd64_v1/overlay rectangle --input=_test/test.jpg --point=650,100 --size=150,125 --colour=#FF0000 --fill --output=dist/overlay_linux_amd64_v1/filled.png
	@OVERLAY_LOG_LEVEL=d dist/overlay_linux_amd64_v1/overlay rectangle --input=_test/test.jpg --point=650,200 --size=150,125 --colour=#00FF00 --stroke=10 --output=dist/overlay_linux_amd64_v1/stroked.png
	@OVERLAY_LOG_LEVEL=d dist/overlay_linux_amd64_v1/overlay rectangle --input=_test/test.jpg --point=650,300 --size=150,125 --colour=#FFFFFF --output=dist/overlay_linux_amd64_v1/default.png
	@OVERLAY_LOG_LEVEL=d dist/overlay_linux_amd64_v1/overlay rectangle --input=_test/test.jpg --point=650,400 --size=150,125 --colour=#FFFFFF --fill --radius=5 --output=dist/overlay_linux_amd64_v1/rounded-filled.png
	@OVERLAY_LOG_LEVEL=d dist/overlay_linux_amd64_v1/overlay rectangle --input=_test/test.jpg --point=650,500 --size=150,125 --colour=#FFFFFF --stroke=10 --radius=5 --output=dist/overlay_linux_amd64_v1/rounded-stroked.png

.PHONY: test-text
test-text: compile # overlay text on top of an image
	@OVERLAY_LOG_LEVEL=d dist/overlay_linux_amd64_v1/overlay text --point=650,100 --size=72 --font=_test/Economica/Economica-Regular.ttf --colour=#FFFFFF --input=_test/test.jpg --output=dist/overlay_linux_amd64_v1/out.png --text="HALLO, WORLD!"

.PHONY: test-image
test-image: compile # overlay an image on top of an image
	@OVERLAY_LOG_LEVEL=d dist/overlay_linux_amd64_v1/overlay image --point=650,100 --input=_test/test.jpg --output=dist/overlay_linux_amd64_v1/out.jpg --image=_test/apple.png

.PHONY: test-pipe
test-pipe: compile # overlay images and text on top of an image
	@cat _test/test.jpg | \
	OVERLAY_LOG_LEVEL=d dist/overlay_linux_amd64_v1/overlay image --point=460,25 --image=_test/apple.png --format=jpg | \
	OVERLAY_LOG_LEVEL=d dist/overlay_linux_amd64_v1/overlay text --point=600,100 --size=72 --font=_test/Economica/Economica-Regular.ttf --colour=#FFFFFF --format=jpg --text="HALLO, WORLD..." | \
	OVERLAY_LOG_LEVEL=d dist/overlay_linux_amd64_v1/overlay text --point=700,160 --size=48 --font=_test/Economica/Economica-Regular.ttf --colour=#00FF0033 --output=dist/overlay_linux_amd64_v1/out.jpg --text="... from me!"
