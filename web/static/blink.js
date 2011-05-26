var blink_delay = 1000;
var elems = document.all.tags("blink");
var state = "visible"

function do_blink() 
{
	if (state == "hidden") {
		state = "visible";
	}	else {
		state = "hidden";
	}

	for(i = 0; i < elems.length; i++) {
		elems[i].style.visibility = state;
	}

	setTimeout("do_blink()", blink_delay);
}

do_blink()

