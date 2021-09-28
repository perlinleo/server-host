const signupFetch = (form, emailInput, passwordInput, repeatPasswordInput) => {
    form.addEventListener('submit', (e) => {
        e.preventDefault();
        const testEmail = emailRegExp.test(emailInput.value);
        const testPassword = passwordRegExp.test(passwordInput.value);
        const testPasswordRepeat = passwordInput.value === repeatPasswordInput.value;

        if (!testEmail) {
            emailInput.className = 'form-field-novalid';
        }

        if (!testPassword) {
            passwordInput.className = 'form-field-novalid';
        }

        if (!testPasswordRepeat) {
            repeatPasswordInput.className = 'form-field-novalid';
        }

        if (!testEmail || !testPassword || !testPasswordRepeat) {
            return;
        }

        const email = emailInput.value.trim();
        const password = passwordInput.value.trim();
        const requestOptions = {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                'email': email,
                'password': password,
            })
        };
        fetch(`${localAddress}/api/v1/signup`, requestOptions)
            .then(response =>
                response.json().then(data => ({
                    data: data,
                    status: response.status
                })).then(res => {
                    console.log(res.data.status)
                    console.log(res.data.body)
                    if (res.status === 200 && res.data.status === 200) {
                        clearRoot();
                        setUserProfile(res.data.body);
                        renderEdit();
                    } else if (res.data.status === 404) {
                        const userNotFound = document.createElement('span')
                        userNotFound.textContent = 'Вы уже зарегестрированы'
                        userNotFound.style.marginTop = "10px"
                        form.appendChild(userNotFound)
                    }
                })).catch((error) => console.log(error));
    })
}