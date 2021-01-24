const Controller = {
    newLine: (text) => {
	const lineDiv = document.createElement("div");
	lineDiv.className = "line";
	lineDiv.textContent = text;
	return lineDiv;
    },

    newMatchedLine: (text) => {
	const lineDiv = document.createElement("div");
	lineDiv.className = "line-matched";
	lineDiv.textContent = text;
	return lineDiv;
    },

    update: (resp) => {
	resp.json().then((resp) => {

	    const resultsDiv = document.querySelector("#results");
	    resultsDiv.textContent = "";
	    const template = document.querySelector("#results-template");


	    resp.results.forEach(res => {
		if (res.matches.length === 0) {
		    return;
		}

		const clone = template.content.cloneNode(true);
		clone.querySelector(".title").textContent = res.title;

		const matches = document.createElement("div");
		matches.className = "matches";


		res.matches.forEach((match) => {
		    const matchLine = match.matchLine;
		    const matchDiv = document.createElement("p");
		    matchDiv.className = "match";

		    match.text.forEach((line, index) => {
			let lineDiv;
			if (index === matchLine) {
			    lineDiv = Controller.newMatchedLine(line);
			} else {
			    lineDiv = Controller.newLine(line);
			}
			matchDiv.appendChild(lineDiv);
		    });

		    matches.appendChild(matchDiv);
		});

		clone.getElementById("result").appendChild(matches);
		resultsDiv.appendChild(clone);
	    });


	});
    },

    search: (ev) => {
	ev.preventDefault();
	const form = document.getElementById("search-form");
	const data = Object.fromEntries(new FormData(form));
	console.log(form, data, data.search_text);

	fetch(`/search?search=${data.query}`).then(Controller.update);
    },
};

const form = document.getElementById("search-form");
form.addEventListener("submit", Controller.search);
