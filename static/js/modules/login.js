const loginFetch = (form, emailInput, passwordInput) => {
    form.addEventListener('submit', (e) => {
        e.preventDefault();
        const testEmail = emailRegExp.test(emailInput.value);
        const testPassword = passwordRegExp.test(passwordInput.value);

        if (!testEmail) {
            emailInput.className = 'form-field-novalid';
        }

        if (!testPassword) {
            passwordInput.className = 'form-field-novalid';
        }
        if (!testEmail || !testPassword) {
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
        fetch(`${localAddress}/api/v1/login`, requestOptions)
            .then(response =>
                response.json().then(data => ({
                    data: data,
                    status: response.status
                })).then(res => {
                    if (res.status === 200 && res.data.status === 200) {
                        clearRoot();
                        setUserProfile(res.data.body);
                        renderFeed();
                        addMenu('feed');
                    } else if (res.data.status === 404) {
                        const userNotFound = document.createElement('span')
                        userNotFound.textContent = 'Вы не зарегестрированы'
                        userNotFound.style.marginTop = "10px"
                        form.appendChild(userNotFound)
                    }
                })).catch((error) => console.log(error));
    })
}