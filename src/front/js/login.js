const signUpButton = document.getElementById('signUp');
const signInButton = document.getElementById('signIn');
const container = document.getElementById('container');

const signInSubmit = document.getElementById('submitSignIn');
const signUpSubmit = document.getElementById('submitSignUp');

const emailSignIn = document.getElementById('email-in');
const passwordSignIn = document.getElementById('password-in');

const firstNameSignUp = document.getElementById('firstName-up');
const lastNameSignUp = document.getElementById('lastName-up');
const emailSignUp = document.getElementById('email-up');
const passwordSignUp = document.getElementById('password-up');

const location1 = document.getElementById('location-1');
const location2 = document.getElementById('location-2');
const location3 = document.getElementById('location-3');

const education1 = document.getElementById('eduction-1');
const education2 = document.getElementById('eduction-2');
const education3 = document.getElementById('eduction-3');




signUpButton.addEventListener('click', () => {
	container.classList.add("right-panel-active");
});

signInButton.addEventListener('click', () => {
	container.classList.remove("right-panel-active");
});


signUpSubmit.addEventListener('click', () => {
	const requestDto = { 
		"username" : document.getElementById('username-up').value,
		"email" : document.getElementById('email-up').value,
		"password": document.getElementById('password-up').value,
		"user" : { 
			"age": parseInt(document.getElementById('age-up').value),
			"race" : document.getElementById('race-up').value,
			"gender": document.getElementById('gender-up').value
		},
		"locations" : [
			{
				"name"     : location1.value.split('-')[0],
				"region"   : location1.value.split('-')[1],
				"district" : location1.value.split('-')[2]
			},
		],
		"education_place": [
			{
				"name": education1.value,
				"education_program": {
					"field": document.getElementById('field-up-1').value,
					"level": document.getElementById('level-up-1').value
				}
				
			}
		]
	};
	if (location2.value != 'na') {
		requestDto.locations.push(
			{
			"name"     : location2.value.split('-')[0],
			"region"   : location2.value.split('-')[1],
			"district" : location2.value.split('-')[2]
		}
		);
	}

	if (location3.value != 'na') {
		requestDto.locations.push(
			{
			"name"     : location3.value.split('-')[0],
			"region"   : location3.value.split('-')[1],
			"district" : location3.value.split('-')[2]
		}
		);
	}

	if (education2.value != 'na') {
		requestDto.education_place.push(
			{
				"name": education2.value,
				"education_program": {
					"field": document.getElementById('field-up-2').value,
					"level": document.getElementById('level-up-2').value
				}
				
			}
		);
	}

	if (education3.value != 'na') {
		requestDto.education_place.push(
			{
				"name": education3.value,
				"education_program": {
					"field": document.getElementById('field-up-3').value,
					"level": document.getElementById('level-up-3').value
				}
				
			}
		);
	}

	const requestEntity = new Request('http://localhost:8080/accounts', {
		mode: 'cors',
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify(requestDto)
	});

	fetch(requestEntity)
		.then(response => {
				return response.json();
		})
		.then(response => {
			if (response.error == null) {
				localStorage.setItem("uuid", response.uuid);
				localStorage.setItem("username", response.username);
				localStorage.setItem("email", response.email);

				var xhr = new XMLHttpRequest();
				xhr.open('POST', 'http://localhost:8080/sessions');
				var requestDto = { 
					"username" : document.getElementById('username-up').value,
					"email" : document.getElementById('email-up').value,
					"password": document.getElementById('password-up').value,
				};
				xhr.credentials = 'same-origin';
				xhr.withCredentials = true;
				xhr.onload = function() {
					if (xhr.status == 200) { 
						window.location.replace("http://localhost/test/task-manager-master/poll.html"); // todo replace
					} 
				  };
				xhr.send(JSON.stringify(requestDto));
			} else {
				window.alert(response.error)
			}
			// updateTasks(response);
			
		}).catch(error => {
			console.error(error);
		});

	

});

signInSubmit.addEventListener('click', () => {
	var xhr = new XMLHttpRequest();
	xhr.open('POST', 'http://localhost:8080/sessions');
	var requestDto = { 
		"username" : document.getElementById('username-in').value,
		"email" : document.getElementById('email-in').value,
		"password": document.getElementById('password-in').value,
	};
	xhr.credentials = 'same-origin';
	xhr.withCredentials = true;
	xhr.onload = function() {
		if (xhr.status == 200) { 
			window.location.replace("http://localhost/test/task-manager-master/poll.html"); // todo replace
		} 
	  };
	xhr.send(JSON.stringify(requestDto));
});

function listLocations() {
	const requestEntity = new Request('http://localhost:8080/locations', {
		mode: 'cors',
		method: 'GET',
		headers: { 'Origin': 'google.com', 'Content-Type': 'application/json' },
	});

	fetch(requestEntity)
		.then(response => {
			return response.json();
		})
		.then(response => {
			setLocations(response);

			for (i = 0; i < response.length; i++) {
				var htmlOption = document.createElement('option');
				htmlOption.textContent = response[i].name + ', ' + response[i].region + ', ' + response[i].district;
				htmlOption.value = response[i].name + '-' + response[i].region + '-' + response[i].district;
				location1.appendChild(htmlOption);
				htmlOption = document.createElement('option');
				htmlOption.textContent = response[i].name + ', ' + response[i].region + ', ' + response[i].district;
				htmlOption.value = response[i].name + '-' + response[i].region + '-' + response[i].district;
				location2.appendChild(htmlOption);
				htmlOption = document.createElement('option');
				htmlOption.textContent = response[i].name + ', ' + response[i].region + ', ' + response[i].district;
				htmlOption.value = response[i].name + '-' + response[i].region + '-' + response[i].district;
				location3.appendChild(htmlOption);
			}
			
		}).catch(error => {
			console.error(error);
		});
}

function setLocations(options) {
	localStorage.setItem("locations", JSON.stringify(options));
}

function getLocations() {
	JSON.parse(localStorage.getItem("locations"));
}

function listEducations() {
	const requestEntity = new Request('http://localhost:8080/education', {
		mode: 'cors',
		// origin: 'google.com',
		method: 'GET',
		headers: { 'Origin': 'google.com', 'Content-Type': 'application/json' },
	});

	fetch(requestEntity)
		.then(response => {
			return response.json();
		})
		.then(response => {
			for (i = 0; i < response.length; i++) {
				var htmlOption = document.createElement('option');
				htmlOption.textContent = response[i].name;
				htmlOption.value = response[i].name;
				education1.appendChild(htmlOption);
				htmlOption = document.createElement('option');
				htmlOption.textContent = response[i].name;
				htmlOption.value = response[i].name;
				education2.appendChild(htmlOption);
				htmlOption = document.createElement('option');
				htmlOption.textContent = response[i].name;
				htmlOption.value = response[i].name;
				education3.appendChild(htmlOption);
			}
			
			setEducations(response);
			// updateView();
		}).catch(error => {
			console.error(error);
		});
}

function setEducations(options) {
	localStorage.setItem("educations", JSON.stringify(options));
}

function getEducations() {
	JSON.parse(localStorage.getItem("educations"));
}

function updateView() {

}

listLocations();
listEducations();
