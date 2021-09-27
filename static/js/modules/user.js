// eslint-disable-next-line no-unused-vars
const user = {
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

function setUserProfile(body) {
  user.id = body.id;
  user.firstName = body.name;
  user.age = body.age;
  user.photoSrc = body.imgSrc;
  user.colorFrom = 'grey';
  user.colorTo = 'black';
  user.text = body.description;
  user.tags = body.tags;
}