const form = document.getElementById('login-form');
const error = document.getElementById('error');

form.addEventListener('submit', (e) => {
  e.preventDefault();

  const username = document.getElementById('username').value;
  const password = document.getElementById('password').value;

  let messages = [];
  if (password.length <= 8) {
    messages.push('Password must be greater than 8 characters');
  } else if (password.length >= 20) {
    messages.push('Password must be less than 20 characters');
  }
  if (messages.length > 0) {
    e.preventDefault();
    error.innerText = messages.join(', ');
  }
  fetch('/login/', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({ username, password })
  })
  .then(response => response.json())
  .then(data => {
    console.log(data)
  })
  .catch(error => {
    console.log(error);
  });
});
