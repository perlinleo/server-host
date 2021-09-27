/* eslint-disable indent */
/* eslint-disable no-unused-vars */
/* eslint-disable prefer-const */
let counter = 0;
let sample = [];

const addProfile = (data) => {
  sample[counter] = Object()
  sample[counter].id=data.id;
  sample[counter].firstName=data.name;
  sample[counter].age=data.age;
  sample[counter].photoSrc=data.imgSrc;
  sample[counter].colorFrom='grey';
  sample[counter].colorTo='black';
  sample[counter].text=data.description;
  sample[counter].tags=data.tags;
}

const swipeUser = (id) => {
  const requestOptions = {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({
      'id': id,
    }),
    credentials: 'include',
  };
  
  fetch(`${serverAddress}/api/v1/nextswipeuser`, requestOptions)
    .then(response =>
      response.json().then(data => ({
        data: data,
        status: response.status
      })).then(res => {
        console.log(res.data.status)
        console.log(res.data.body)
        if (res.data.status === 200) {
          clearRoot();
          addProfile(res.data.body)
          renderFeed();
          addMenu('feed');
        } else if (res.data.status === 404) {
          clearRoot();
          renderFeed();
          addMenu('feed');
        }
        
        // console.log(sample)
       
        // if (res.status === 200 && res.data.status === 200) {
          
        // }
        // console.log(res.data);
        // if (res.data.status === 'ok') {
        //     profilePage();
        // }
        // console.log(res.data.status)
      })).catch((error) => console.log(error));

}
