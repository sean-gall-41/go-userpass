const username = document.getElementById('username');
const pass = document.getElementById('password');
const form = document.getElementById('login-form');
const error = document.getElementById('error');

form.addEventListener('submit', (e) => {
  let messages = [];
  if (pass.value.length <= 8) {
    messages.push('Password must be greater than 8 characters');
  } else if (pass.value.length >= 20) {
    messages.push('Password must be less than 20 characters');
  }
  if (messages.length > 0) {
    e.preventDefault();
    error.innerText = messages.join(', ');
  }
  // TODO: create AJAX request, handle response Promise and update UI
})
