// eslint-disable-next-line no-unused-vars
const serverAddress = 'http://95.84.192.140'
// const serverAddress = 'http://127.0.0.1'

const user = {
  'id': 1,
  'firstName': 'Malenkaya Yachta',
  'age': '24',
  'photoSrc': 'img/Yachty-tout.jpg',
  'colorFrom': 'grey',
  'colorTo': 'black',
  'text': `
      Miles Parks McCollum, known professionally as 
      Lil Yachty, is an American rapper from Atlanta, Georgia. 
      McCollum first gained recognition on the internet in 2015
      for his singles 'One Night' and 'Minnesota' from his 
      debut EP Summer Songs.
      He released his debut mixtape Lil Boat in March 2016.
      `,
  'tags': [
    'malenkaya',
    'yachta',
    'straight',
    'surfing',
  ],
};

function setUserProfile(data){
  user.id=data.id;
  user.firstName=data.name;
  user.age = data.age;
  user.text = data.description;
  user.photoSrc = data.imgSrc
  user.tags = data.tags
}