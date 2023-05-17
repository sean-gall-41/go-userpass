const username = document.getElementById('name');
const pass = document.getElementById('pass');
const form = document.getElementById('form');
const error = document.getElementById('error');

form.addEventListener('submit', (e) => {
  let messages = [];
  if (username.value === '' || username.value == null) {
    messages.push('Username is required');
  }
  if (pass.value.length <= 8) {
    messages.push('Password must be greater than 8 characters');
  } else if (pass.value.length >= 20) {
    messages.push('Password must be less than 20 characters');
  }
  if (messages.length > 0) {
    e.preventDefault();
    error.innerText = message.join(', ');
  }
})
