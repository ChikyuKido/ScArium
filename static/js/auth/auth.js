
document.getElementById('authForm').addEventListener('submit', function(event) {
    event.preventDefault();
    const username = document.getElementById('username').value;
    const password = document.getElementById('password').value;
    const errorMessageDiv = document.getElementById('error-message');
    if (errorMessageDiv) {
        errorMessageDiv.remove();
    }

    if(window.location.toString().includes("register")) {
        register(username,password)
    }else  if(window.location.toString().includes("login")) {
        login(username,password)
    }else  if(window.location.toString().includes("adminRegister")) {
        adminRegister(username,password)
    }

});

function register(username,password) {
    fetch('/api/v1/auth/register', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username, password })
    })
        .then(response => {
            if (response.status === 201) {
                window.location.href = '/auth/login';
            } else {
                return response.json().then(data => {
                    throw new Error(data.error || 'Registration failed');
                });
            }
        })
        .catch(error => {
            const form = document.getElementById('authForm');
            const errorMessage = document.createElement('p');
            errorMessage.id = 'error-message';
            errorMessage.className = 'help is-danger';
            errorMessage.textContent = error.message || 'An error occurred during registration.';
            form.appendChild(errorMessage);
        });
}
function login(username,password) {
    fetch('/api/v1/auth/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username, password })
    })
        .then(response => {
            if (response.status === 200) {
                window.location.href = '/';
            } else {
                return response.json().then(data => {
                    throw new Error(data.error || 'Registration failed');
                });
            }
        })
        .catch(error => {
            const form = document.getElementById('authForm');
            const errorMessage = document.createElement('p');
            errorMessage.id = 'error-message';
            errorMessage.className = 'help is-danger';
            errorMessage.textContent = error.message || 'An error occurred during registration.';
            form.appendChild(errorMessage);
        });
}
function adminRegister(username,password) {
    fetch('/api/v1/auth/adminRegister', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username, password })
    })
        .then(response => {
            if (response.status === 201) {
                window.location.href = '/auth/login';
            } else {
                return response.json().then(data => {
                    throw new Error(data.error || 'Registration failed');
                });
            }
        })
        .catch(error => {
            const form = document.getElementById('authForm');
            const errorMessage = document.createElement('p');
            errorMessage.id = 'error-message';
            errorMessage.className = 'help is-danger';
            errorMessage.textContent = error.message || 'An error occurred during registration.';
            form.appendChild(errorMessage);
        });
}