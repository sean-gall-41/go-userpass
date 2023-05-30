const loginForm = document.getElementById('login-form');
const registerForm = document.getElementById('register-form');
const usernameForgetForm = document.getElementById('username-forget-form');
const passwordResetRequestForm = document.getElementById('password-reset-request-form');
const passwordResetForm = document.getElementById('password-reset-form');
const error = document.getElementById('error');

if (loginForm !== null ) {
  loginForm.addEventListener('submit', (e) => {
    e.preventDefault();
    const username = document.getElementById('username').value;
    const password = document.getElementById('password').value;
    fetch('/login/', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ username, password })
    })
    .then(response => response.json())
    .then(data => {
      if (data.success) {
        window.location.href = '/login-success/';
      } else {
        const errorElement = document.getElementById("error");
        errorElement.innerText = data.message;
        errorElement.style.display = 'block';
      }
    })
    .catch(error => {
      console.log(error);
    });
  });
}

if (registerForm !== null) {
  registerForm.addEventListener('submit', (e) => {
    e.preventDefault();
    const email = document.getElementById('email').value;
    const username = document.getElementById('username').value;
    const password = document.getElementById('password').value;
    fetch('/register/', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ email, username, password })
    })
    .then(response => response.json())
    .then(data => {
      if (data.success) {
        window.location.href = '/register-success/';
      } else {
        const errorElement = document.getElementById("error");
        errorElement.innerText = data.message;
        errorElement.style.display = 'block';
      }
    })
    .catch(error => {
      console.log(error);
    });
  });
}

if (usernameForgetForm !== null) {
  usernameForgetForm.addEventListener('submit', (e) => {
    e.preventDefault();
    const email = document.getElementById('email').value;
    fetch('/username-forget/', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ email })
    })
    .then(response => response.json())
    .then(data => {
      if (data.success) {
        window.location.href = '/username-forget-success/';
      } else {
        const errorElement = document.getElementById("error");
        errorElement.innerText = data.message;
        errorElement.style.display = 'block';
      }
    })
    .catch(error => {
      console.log(error);
    });
  });
}

if (passwordResetRequestForm !== null) {
  passwordResetRequestForm.addEventListener('submit', (e) => {
    e.preventDefault();
    const email = document.getElementById('email').value;
    fetch('/password-reset-request/', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ email })
    })
    .then(response => response.json())
    .then(data => {
      if (data.success) {
        window.location.href = '/password-reset-request-success/';
      } else {
        const errorElement = document.getElementById("error");
        errorElement.innerText = data.message;
        errorElement.style.display = 'block';
      }
    })
    .catch(error => {
      console.log(error);
    });
  });
}

if (passwordResetForm !== null) {
  passwordResetForm.addEventListener('submit', (e) => {
    e.preventDefault();
    const password = document.getElementById('password').value;
    // TODO: handle password-verify not being equal to password
    fetch(window.location.href, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ password })
    })
    .then(response => response.json())
    .then(data => {
      if (data.success) {
        // TODO: redirect to success screen
        console.log("Successfully reset password.")
        //window.location.href = '/login-success/';
      } else {
        const errorElement = document.getElementById("error");
        errorElement.innerText = data.message;
        errorElement.style.display = 'block';
      }
    })
    .catch(error => {
      console.log(error);
    });
  });
}

