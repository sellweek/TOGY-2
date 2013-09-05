"use strict"

var data = {
    timestamp: "2013-08-20 18:43",
    sections: [{
        title: "Prvá sekcia",
        showTitle: true,
        slides: [{
            title: "Slajd jeden",
            contents: "<ul><li>Jedna</li><li>Dva</li><li>Tri</li></ul>"
        },
        {
            title: "Slajd dva",
            contents: "<ul><li>Tri</li><li>Štyri</li><li>Päť</li></ul>"
        },
        {
            title: "Slajd s obrázkom",
            contents: 'Vario LF počas homologizácie<br><img src="http://imhd.zoznam.sk/ke/media/gn/00074389/701-Ludvikov-dvor-rychlodraha.jpg">'
        }],
        start: "2013-08-20 18:20",
        //end: "2013-08-20 19:00",
        delays: 1
    },
    {
        title: "Druhá sekcia",
        showTitle: false,
        slides: [{
            title: "Začiatok druhej sekcie",
            contents: "<h3>Bum!</h3>"
        },
        {
            title: "Skúška tabuľky",
            contents: '<table border="1"><tr><td>ahooj</td><td>ahooj</td></tr><tr><td>ahooj</td><td>ahooj</td></tr></table>'
        },
        {
            title: "",
            contents: "<h1>Zbohom</h1>"
        }],
        start: "2013-08-19 18:20",
        delays: 3
    }]
};

var TOGY = function() {

    //timeoutID is the ID of timeout used by cycle()
    var timeoutID = undefined;

    //Coords encapsulates Reveal.js presentation indices.
    //Every slide has coords, using which it can be shown by
    //Reveal.slide() function.
    //Coords of a section are those of its first slide.
    //In current implementation, f (fragment index) 
    //will always be undefined.
    var Coords = function(h, v) {
        this.h = h;
        this.v = v;
        this.f = undefined;
    };

    //isLastInSection returns whether given Coords belong to
    //the last slide in a section.
    Coords.prototype.isLastInSection = function(sections) {
        return this.v == (sectionLength(sections[this.h]) - 1);
    };

    //show instructs Reveal to move to slide located at given Coords.
    Coords.prototype.show = function() {
    	Reveal.slide(this.h, this.v, undefined);
    }

    //getCurrent returns Coords of currently active slide, using Reveal.getIndices().
    Coords.getCurrent = function() {
        var ind = Reveal.getIndices();
        return new Coords(ind.h, ind.v);
    }

    //sectionLength returns number of slides in a section.
    //This function is necessary, because setting showTitle
    //to true adds one slide (which doesn't appear in the slides array)
    //to the section.
    var sectionLength = function(section) {
    	var slides = section.slides.length;
    	if (section.showTitle) {
    		slides += 1;
    	}
    	return slides;
    };

    //showNext shows next slide and returns its Coords.
    var showNext = function(sections) {
        var c = Coords.getCurrent();
        if (!c.isLastInSection(data.sections)) {
            c.v += 1;
            c.show();
            return c;
        }

        var nextCoords = nextActiveSection(c, sections);
        nextCoords.show();
        return nextCoords;
    };

    //cycle cycles between slides by calling showNext in the
    //interval given by respective sections.
    var cycle = function(sections) {
    	var c = showNext(sections);
    	var timeout = sections[c.h].delays * 1000;
    	timeoutID = setTimeout(cycle, timeout, sections);
    };

    //stop stops cycling between intervals.
    var stop = function() {
        if (timeoutID) {
            window.clearTimeout(timeoutID);
            timeoutID = undefined;
            return true;
        }
        return false;
    }

    //nextActiveSection returns coordinates of the next active section.
    //If no section is active, it returns coordinates of the first slide
    //of the current section.
    //A section is considered active if current time is between its
    //start and end times. 
    var nextActiveSection = function(currCoords, sections) {
    	var now = Date.now();
    	for (var i = currCoords.h + 1; i < sections.length; i++) {
    		if (isActiveAt(now, sections[i])) {
    			return new Coords(i, 0);
    		}
    	}
    	//If we can't find an active section before the end of broadcast, we have
        //to start from the beginning.
    	for (i = 0; i <= currCoords.h; i++) {
    		if (isActiveAt(now, sections[i])) {
    			return new Coords(i, 0);
    		}
    	}

    	return new Coords(currCoords.h, 0)
    };

    //buildHTML constructs DOM used by Reveal.js into
    //a div, class of which is "slides"
    var buildHTML = function(sections) {
        var slidesArea = $(".slides");
        slidesArea.empty();

        for (var i = 0; i < sections.length; i++) {
            var container = document.createElement("section");

            if (sections[i].showTitle) {
                var titleSection = document.createElement("section");
                $(titleSection).append("<h1>" + sections[i].title + "</h1>");
                $(container).append(titleSection);
            }

            var slides = sections[i].slides;
            for (var j = 0; j < slides.length; j++) {
                var slide = document.createElement("section");
                $(slide).append("<h2>" + slides[j].title + "</h2>").append(slides[j].contents);
                $(container).append(slide);
            }

            $(slidesArea).append(container);
        }
    };

    //buildCoords adds coords to slides and sections in passed in object.
    var buildCoords = function(sections) {
        for (var i = 0; i < sections.length; i++) {
            sections[i].coords = new Coords(i, 0)

            var slides = sections[i].slides;
            for (var j = 0; j < slides.length; j++) {
                var v;
                if (sections[i].showTitle) {
                    v = j + 1;
                } else {
                    v = j;
                }
                slides[j].coords = new Coords(i, v);
            }
        }
    };

    //replaceDates replaces string encoded dates and times
    //(in YYYY-MM-DD HH:MM) format in passed in object
    //by JavaScript Date objects.
    var replaceDates = function(data) {
        data.timestamp = new Date(data.timestamp)
        for (var i = 0; i < data.sections.length; i++) {
            if (data.sections[i].start) {
                data.sections[i].start = new Date(data.sections[i].start);
            }
            if (data.sections[i].end) {
                data.sections[i].end = new Date(data.sections[i].end);
            }
        }
    };

    //isActiveAt returns whether a section is
    //active on the given Date.
    var isActiveAt = function(date, section) {
        if (section.start && date < section.start) {
            return false;
        }

        if (section.end && date > section.end) {
            return false;
        }

        return true;
    };

    //initialize prepares data object and DOM
    //and starts cycling the slides.
    var initialize = function(data) {
    	replaceDates(data);
    	buildCoords(data.sections);
    	buildHTML(data.sections);
    	Reveal.initialize({
		    controls: false,
		    progress: false,
		    history: false,
		    keyboard: false,
		    touch: false,
		    overview: false,
		    center: true,
		    loop: false,
		    rtl: false,
		    autoSlide: 0,
		    mouseWheel: false,
		    transition: 'linear',
		   	transitionSpeed: 'default',
		    backgroundTransition: 'default'
		});
		cycle(data.sections);
    };

    return {
        "initialize": initialize,
        "stop": stop
    }
}();