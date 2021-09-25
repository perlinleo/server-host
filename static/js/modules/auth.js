const serverAddress = 'http://95.84.192.140'
const localAddress = 'http://127.0.0.1'


const loginWithCookie = () => {
  const requestOptions = {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json'
    },
  };
  fetch(`${localAddress}/api/v1/cookie`, requestOptions)
    .then(response =>
      response.json().then(data => ({
        data: data,
        status: response.status
      })).then(res => {
        if (res.status === 200 && res.data.status === 200) {
          clearRoot();

          setUserProfile(res.data.body);
          userProfileRender();
          addMenu('profile');
        }
        console.log(res.data);
        // if (res.data.status === 'ok') {
        //     profilePage();
        // }
        // console.log(res.data.status)
      })).catch((error) => console.log(error));

}

function loginWithCredentials(email, password){
    const requestOptions = {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          'email': email,
          'password': password,
        })
      };
      fetch(`${serverAddress}/api/v1/login`, requestOptions)
        .then(response =>
          response.json().then(data => ({
            data: data,
            status: response.status
          })).then(res => {
            if (res.status === 200 && res.data.status === 200) {
              loginWithCookie();
            } else if (res.data.status === 404) {
              loginPageError("User not found")
            }
          })).catch((error) => console.log(error));
}



function logoutCookie(){
    
    const requestOptions = {
        method: 'GET',
        headers: { 'Content-Type': 'application/json' },
      };
      fetch(`${serverAddress}/api/v1/logout`, requestOptions)
          clearRoot();
         
}