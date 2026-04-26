.PHONY: test-draw-canvas
test-draw-canvas: compile # create a canvas with the given size and colour
	@OVERLAY_LOG_LEVEL=d dist/overlay_linux_amd64_v1/overlay draw canvas --size=640,480 --colour=#FF0000 --output=dist/overlay_linux_amd64_v1/red.png
	@OVERLAY_LOG_LEVEL=d dist/overlay_linux_amd64_v1/overlay draw canvas --size=640,480 --colour=#00FF00 --output=dist/overlay_linux_amd64_v1/green.jpg
	@OVERLAY_LOG_LEVEL=d dist/overlay_linux_amd64_v1/overlay draw canvas --size=640,480 --colour=#0000FF --output=dist/overlay_linux_amd64_v1/blue.bmp

.PHONY: test-draw-rectangle
test-draw-rectangle: compile # create a rectangle with the given size and colour
	@OVERLAY_LOG_LEVEL=d dist/overlay_linux_amd64_v1/overlay draw rectangle --input=_test/test.jpg --point=650,100 --size=150,125 --colour=#FF0000 --fill --output=dist/overlay_linux_amd64_v1/filled.png
	@OVERLAY_LOG_LEVEL=d dist/overlay_linux_amd64_v1/overlay draw rectangle --input=_test/test.jpg --point=650,200 --size=150,125 --colour=#00FF00 --stroke=10 --output=dist/overlay_linux_amd64_v1/stroked.png
	@OVERLAY_LOG_LEVEL=d dist/overlay_linux_amd64_v1/overlay draw rectangle --input=_test/test.jpg --point=650,300 --size=150,125 --colour=#FFFFFF --output=dist/overlay_linux_amd64_v1/default.png
	@OVERLAY_LOG_LEVEL=d dist/overlay_linux_amd64_v1/overlay draw rectangle --input=_test/test.jpg --point=650,400 --size=150,125 --colour=#FFFFFF --fill --radius=5 --output=dist/overlay_linux_amd64_v1/rounded-filled.png
	@OVERLAY_LOG_LEVEL=d dist/overlay_linux_amd64_v1/overlay draw rectangle --input=_test/test.jpg --point=650,500 --size=150,125 --colour=#FFFFFF --stroke=2 --radius=15 --output=dist/overlay_linux_amd64_v1/rounded-stroked.png

.PHONY: test-draw-circle
test-draw-circle: compile # create a circle with the given size and colour
	@OVERLAY_LOG_LEVEL=d dist/overlay_linux_amd64_v1/overlay draw circle --input=_test/test.jpg --point=650,200 --colour=#FF0000 --fill --radius=100 --output=dist/overlay_linux_amd64_v1/filled.png
	@OVERLAY_LOG_LEVEL=d dist/overlay_linux_amd64_v1/overlay draw circle --input=_test/test.jpg --point=650,200 --colour=#00FF00 --stroke=10 --radius=100 --output=dist/overlay_linux_amd64_v1/stroked.png

.PHONY: test-draw-circular-arc
test-draw-circular-arc: compile # create a circular arc with the given size and colour
	@OVERLAY_LOG_LEVEL=d dist/overlay_linux_amd64_v1/overlay draw circular-arc --input=_test/test.jpg --point=650,200 --colour=#FF0000 --fill --radius=100 --angle=0,90 --output=dist/overlay_linux_amd64_v1/circular-arc-90-degrees.png
	@OVERLAY_LOG_LEVEL=d dist/overlay_linux_amd64_v1/overlay draw circular-arc --input=_test/test.jpg --point=650,200 --colour=#00FF00 --stroke=10 --radius=100 --angle=0,270 --output=dist/overlay_linux_amd64_v1/circular-arc-270-degrees.png

.PHONY: test-draw-elliptical-arc
test-draw-elliptical-arc: compile # create an elliptical arc with the given size and colour
	@OVERLAY_LOG_LEVEL=d dist/overlay_linux_amd64_v1/overlay draw elliptical-arc --input=_test/test.jpg --point=650,200 --colour=#FF0000 --fill --radius=100,50 --angle=0,90 --output=dist/overlay_linux_amd64_v1/elliptical-arc-90-degrees.png
	@OVERLAY_LOG_LEVEL=d dist/overlay_linux_amd64_v1/overlay draw elliptical-arc --input=_test/test.jpg --point=650,200 --colour=#00FF00 --stroke=10 --radius=100,50 --angle=0,270 --output=dist/overlay_linux_amd64_v1/elliptical-arc-270-degrees.png

.PHONY: test-draw-ellipse
test-draw-ellipse: compile # create an ellipse with the given size and colour
	@OVERLAY_LOG_LEVEL=d dist/overlay_linux_amd64_v1/overlay draw ellipse --input=_test/test.jpg --point=650,200 --colour=#FF0000 --fill --radius=100,50 --output=dist/overlay_linux_amd64_v1/filled.png
	@OVERLAY_LOG_LEVEL=d dist/overlay_linux_amd64_v1/overlay draw ellipse --input=_test/test.jpg --point=650,200 --colour=#00FF00 --stroke=10 --radius=100,50 --output=dist/overlay_linux_amd64_v1/stroked.png

.PHONY: test-draw-text
test-draw-text: compile # overlay text on top of an image
	@OVERLAY_LOG_LEVEL=d dist/overlay_linux_amd64_v1/overlay draw text --point=650,100 --size=72 --font=_test/Economica/Economica-Regular.ttf --colour=#FFFFFF --input=_test/test.jpg --output=dist/overlay_linux_amd64_v1/out.png --text="HALLO, WORLD!"

.PHONY: test-draw-image
test-draw-image: compile # overlay an image on top of an image
	@OVERLAY_LOG_LEVEL=d dist/overlay_linux_amd64_v1/overlay draw image --point=650,100 --input=_test/test.jpg --output=dist/overlay_linux_amd64_v1/out.jpg --image=_test/apple.png

.PHONY: test-draw-pipeline
test-draw-pipeline: compile # overlay images and text on top of an image
	@cat _test/test.jpg | \
	OVERLAY_LOG_LEVEL=d dist/overlay_linux_amd64_v1/overlay draw image --point=460,25 --image=_test/apple.png --format=jpg | \
	OVERLAY_LOG_LEVEL=d dist/overlay_linux_amd64_v1/overlay draw text --point=600,100 --size=72 --font=_test/Economica/Economica-Regular.ttf --colour=#FFFFFF --format=jpg --text="HALLO, WORLD..." | \
	OVERLAY_LOG_LEVEL=d dist/overlay_linux_amd64_v1/overlay draw text --point=700,160 --size=48 --font=_test/Economica/Economica-Regular.ttf --colour=#00FF0033 --output=dist/overlay_linux_amd64_v1/out.jpg --text="... from me!"

.PHONY: test-transform-rotate
test-transform-rotate: compile # rotate the image
	@OVERLAY_LOG_LEVEL=d dist/overlay_linux_amd64_v1/overlay transform rotate --input=_test/test.jpg --angle=45 --pivot=10,10 --resize --output=dist/overlay_linux_amd64_v1/rotated-res.png

.PHONY: test-transform-fliph
test-transform-fliph: compile # flip an image horizontally
	@OVERLAY_LOG_LEVEL=d dist/overlay_linux_amd64_v1/overlay transform fliph --input=_test/test.jpg --output=dist/overlay_linux_amd64_v1/flipped-horizontally.png

.PHONY: test-transform-flipv
test-transform-flipv: compile # flip an image vertically
	@OVERLAY_LOG_LEVEL=d dist/overlay_linux_amd64_v1/overlay transform flipv --input=_test/test.jpg --output=dist/overlay_linux_amd64_v1/flipped-vertically.png

.PHONY: test-transform-zoom
test-transform-zoom: compile # zoom the image
	@OVERLAY_LOG_LEVEL=d dist/overlay_linux_amd64_v1/overlay transform zoom --input=_test/test.jpg --factor=2.0 --pivot=10,10 --output=dist/overlay_linux_amd64_v1/zoomed-in.png
	@OVERLAY_LOG_LEVEL=d dist/overlay_linux_amd64_v1/overlay transform zoom --input=_test/test.jpg --factor=0.5 --pivot=10,10 --output=dist/overlay_linux_amd64_v1/zoomed-out.png
