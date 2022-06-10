function updateView(options) {
  updateOptions(options);
  if (options != null) {
    for (let i = 0; i < options.length; i++) {
      let option = options[i];
        $("#option-list-analysis").append(`
                <li class="option-item">${option}
                <select class="rating" id="rate-${i}">
                    <option value="10">10</option>
                    <option value="9">9</option>
                    <option value="8">8</option>
                    <option value="7">7</option>
                    <option value="6">6</option>
                    <option value="5">5</option>
                    <option value="4">4</option>
                    <option value="3">3</option>
                    <option value="2">2</option>
                    <option value="1">1</option>
                </select>
                </li>
            `);
    }
  }
}

$("#poll-submit").on("click", submitOption);

function submitOption() {
  options = getOptions();
  const requestDto = {
    'poll' : []

  };
  for (let i = 0; i < options.length; i++) {
    let option = options[i];
    let value = $("#rate-" + i).val();
    requestDto.poll.push(
      {
        "word": option,
        "mark": value
      }
    );
  }

  var xhr = new XMLHttpRequest();
	xhr.open('POST', 'http://localhost:8080/authorized/poll');
	xhr.credentials = 'same-origin';
	xhr.withCredentials = true;

	xhr.onload = function() {
    console.log(xhr.status);

    if (xhr.status == 201) { 

		} 
	  };

	xhr.send(JSON.stringify(requestDto));
}



function listOptions() {
  var xhr = new XMLHttpRequest();
	xhr.open('GET', 'http://localhost:8080/authorized/poll');
	xhr.credentials = 'same-origin';
	xhr.withCredentials = true;

	xhr.onload = function() {
		if (xhr.status == 201) { 
      updateView(JSON.parse(xhr.response));
		} 
	  };

	xhr.send();
}

function updateOptions(options) {
  localStorage.setItem("options", JSON.stringify(options));
}

function getOptions() {
  return JSON.parse(localStorage.getItem("options"));
}

// listDashboards();
listOptions();
// localStorage.setItem("user-uuid", "PASTE HERE THE UUID");


// loadDataSelect();

// ###### Main #####
