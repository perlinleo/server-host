import LoginComponent from './components/Login/Login.js';

import FeedComponent from './components/Feed/Feed.js';

import ProfileComponent from './components/Profile/Profile.js';

const root = document.getElementById('root');

window.addEventListener('load', (e) => {
  e.preventDefault();
  window.User.loginWithCookie(()=>{
    feedPage();
    addMenu('feed');
  })
})

document.addEventListener('click', clickButtons, false);


loginPage();

const configApp = {
  login: {
    href: '/login',
    name: 'Авторизация',
    open: loginPage,
  },
  signup: {
    href: '/signup',
    name: 'Регистрация',
    open: signupPage,
  },
  'menu-feed': {
    link: '/feed',
    name: 'feed',
    open: feedPage,
  },
  'menu-profile': {
    link: '/profile',
    name: 'profile',
    open: profilePage,
  },
  'expand-card': {
    link: '/feed',
    name: 'feed',
    open: profileRender,
  },
  'shrink-card': {
    link: '/feed',
    name: 'feed',
    open: feedPage,
  },
  'menu-likes': {
    link: '/likes',
    name: 'likes',
    open: notDoneYet,
  },
  'menu-chat': {
    link: '/likes',
    name: 'chat',
    open: notDoneYet,
  },
  'profile-edit': {
    link: '/profile/edit',
    name: 'edit profile',
    open: notDoneYet,
  },
  'profile-logout': {
    link: '/',
    name: 'logging out',
    open: logout,
  }
}

function feedPage() {
  const feed = new FeedComponent(root);
  feed.render();
}

function profilePage() {
  root.innerHTML = '';
  const profile = new ProfileComponent(root);
  profile.render();
}

function logout(){
  window.User.logoutCookie(
    ()=>{
      loginPage();
    })
}



function loginPage() {
  root.innerHTML = '';
  const login = new LoginComponent(root);
  login.render();
}

function signupPage() {
  root.innerHTML = '';

  const header = createCenterContainer();

  const headerText = document.createElement('span');
  headerText.textContent = 'Регистрация';
  headerText.classList.add('login-header');

  header.appendChild(headerText);
  root.appendChild(header);

  const form = document.createElement('form');
  form.classList.add('login-form');

  const emailInput = createInput('email', 'Почта', 'email');
  emailInput.className = 'form-field-valid';
  emailInput.addEventListener('input', () => {
    const test = emailInput.value.length === 0 || emailRegExp.test(emailInput.value);
    if (test) {
      emailInput.className = 'form-field-valid';
    } else {
      emailInput.className = 'form-field-novalid'
    }
  })
  const passwordInput = createInput('password', 'Пароль', 'password');
  passwordInput.className = 'form-field-valid';
  passwordInput.addEventListener('input', () => {
    const test = passwordInput.value.length === 0 || passwordRegExp.test(passwordInput.value);

    if (test) {
      passwordInput.className = 'form-field-valid';
    } else {
      passwordInput.className = 'form-field-novalid'
    }
  })

  const repeatPasswordInput = createInput('password', 'Пароль', 'password');
  repeatPasswordInput.className = 'form-field-valid';
  repeatPasswordInput.addEventListener('input', () => {
    const test = passwordInput.value === repeatPasswordInput.value;
    if (test) {
      repeatPasswordInput.className = 'form-field-valid';
    } else {
      repeatPasswordInput.className = 'form-field-novalid'
    }
  })

  // кнопка зарегестрироваться
  const submitButton = document.createElement('button');
  submitButton.type = 'submit';
  submitButton.classList.add('login-button');

  const buttonFilling = createCenterContainer();
  const buttonText = document.createElement('span');
  buttonText.textContent = 'Зарегистрироваться';
  buttonText.classList.add('signup-button-text');
  const buttonIcon = document.createElement('img');
  buttonIcon.src = './svg/next.svg';
  buttonIcon.classList.add('svg-next');

  buttonFilling.appendChild(buttonText);
  buttonFilling.appendChild(buttonIcon);

  submitButton.appendChild(buttonFilling);

  const emailIcon = document.createElement('img');
  emailIcon.src = './svg/email.svg'
  emailIcon.classList.add('input-icon');
  const passwordIcon = document.createElement('img');
  passwordIcon.src = './svg/password.svg';
  passwordIcon.classList.add('input-icon');
  const repeatPasswordIcon = document.createElement('img');
  repeatPasswordIcon.src = './svg/password.svg';
  repeatPasswordIcon.classList.add('input-icon');

  const emailFieldWithIcon = document.createElement('div');
  emailFieldWithIcon.classList.add('input-with-icon');
  const passwordFieldWithIcon = document.createElement('div');
  passwordFieldWithIcon.classList.add('input-with-icon');
  const repeatPasswordFieldWithIcon = document.createElement('div');
  repeatPasswordFieldWithIcon.classList.add('input-with-icon');

  const logoBg = document.createElement('div');
  logoBg.classList.add('drip-logo-bg');

  const formContainer = createCenterContainer();

  emailFieldWithIcon.appendChild(emailInput);
  emailFieldWithIcon.appendChild(emailIcon);
  passwordFieldWithIcon.appendChild(passwordInput);
  passwordFieldWithIcon.appendChild(passwordIcon);
  repeatPasswordFieldWithIcon.appendChild(repeatPasswordInput);
  repeatPasswordFieldWithIcon.appendChild(repeatPasswordIcon);

  logoBg.appendChild(emailFieldWithIcon);
  logoBg.appendChild(passwordFieldWithIcon);
  logoBg.appendChild(repeatPasswordFieldWithIcon);

  form.appendChild(logoBg);
  form.appendChild(submitButton);

  formContainer.appendChild(form);

  const regLink = document.createElement('a');
  regLink.href = '/createProfile';
  regLink.textContent = 'Зарегестрироваться';
  regLink.dataset.section = 'createProfile';

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
    const passwordRepeat = repeatPasswordInput.value.trim();
    const requestOptions = {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        'email': email,
        'password': password,
        'passwordRepeat': passwordRepeat,
      })
    };
    fetch(`${serverAddress}/api/v1/signup`, requestOptions)
      .then(response =>
        response.json().then(data => ({
          data: data,
          status: response.status
        })).then(res => {
          if (res.status === 200 && res.data.status === 200) {
            
            loginPage();
          } else if (res.data.status === 404) {
            const userNotFound = document.createElement('span')
            userNotFound.textContent = 'Вы уже зарегестрированы'
            userNotFound.style.marginTop = "10px"
            form.appendChild(userNotFound)
          }
        })).catch((error) => console.log(error));
  })

  root.appendChild(formContainer);
}

function createProfilePage() {
  root.innerHTML = '<span>Create Profile Page</span>';
}




root.addEventListener('click', (e) => {
  const {
    target
  } = e;
  if (configApp[target.className]) {
   
    configApp[target.className].open();
    addMenu(configApp[target.className].name);
  }
  if (target instanceof HTMLAnchorElement) {
    e.preventDefault();

    configApp[target.dataset.section].open();
  }
})


/**
    Обрабатывает нажатия
 * @param {event} event - Событие
 */
function clickButtons(event) {
  const {
    target
  } = event;
  if (configApp[target.className]) {
   
    configApp[target.className].open();
    addMenu(configApp[target.className].name);
  }
}


/**
 * Функция, которая рендерит несделанные страницы
 */
function notDoneYet() {
  root.innerHTML = 'Not done Yet';
}




// START


/**
 * в будущем
 */
function fillCardMain() {
  const cardMain = document.getElementById('cardMainID');

  const lol = sample[counter];

  const img = document.createElement('img');
  img.src = `${lol.photoSrc}`;
  img.className = 'card-el profile-image-expand';
  cardMain.appendChild(img);

  const divName = document.createElement('div');
  divName.className = 'name';
  divName.textContent = `${lol.firstName}, ${lol.age}`;
  cardMain.appendChild(divName);

  const divBord = document.createElement('div');
  divBord.className = 'card-el bord';
  cardMain.appendChild(divBord);

  const divDesc = document.createElement('div');
  divDesc.className = 'card-el desc';
  divDesc.innerHTML = `${lol.text}`;
  cardMain.appendChild(divDesc);

  const divBord2 = document.createElement('div');
  divBord2.className = 'card-el bord';
  cardMain.appendChild(divBord2);

  const divTags = document.createElement('div');
  divTags.id = 'tagsID';
  divTags.className = 'card-el tags-container';
  cardMain.appendChild(divTags);

  for (const tag in lol.tags) {
    if (Object.prototype.hasOwnProperty.call(lol.tags, tag)) {
      const buttonTag = document.createElement('div');
      buttonTag.className = 'tag';
      buttonTag.innerHTML = `${lol.tags[tag]}`;
      divTags.appendChild(buttonTag);
    }
  }
}

/**
 * ОДНАЖДЫ
 */
function profileRender() {
  const divCrad = document.createElement('div');
  divCrad.id = 'cardID';
  divCrad.className = 'card-expand';
  root.appendChild(divCrad);

  const divTapBar = document.createElement('div');
  divTapBar.id = 'tapbarID';
  divTapBar.className = 'tapbar-container';
  root.appendChild(divTapBar);

  fillCard();
}

/**
 * здесь найдется джсдокс для каждого
 */
function fillUserProfile() {
  const cardMain = document.getElementById('cardMainID');

  const lol = user;

  const img = document.createElement('img');
  img.src = `${lol.photoSrc}`;
  img.className = 'card-el profile-image-expand';
  cardMain.appendChild(img);

  const divName = document.createElement('div');
  divName.className = 'name';
  divName.textContent = `${lol.firstName}, ${lol.age}`;
  cardMain.appendChild(divName);

  const divBord = document.createElement('div');
  divBord.className = 'card-el bord';
  cardMain.appendChild(divBord);

  const divDesc = document.createElement('div');
  divDesc.className = 'card-el desc';
  divDesc.innerHTML = `${lol.text}`;
  cardMain.appendChild(divDesc);

  const divBord2 = document.createElement('div');
  divBord2.className = 'card-el bord';
  cardMain.appendChild(divBord2);

  const divTags = document.createElement('div');
  divTags.id = 'tagsID';
  divTags.className = 'card-el tags-container';
  cardMain.appendChild(divTags);

  for (const tag in lol.tags) {
    if (Object.hasOwnProperty.call(lol.tags, tag)) {
      const buttonTag = document.createElement('button');
      buttonTag.className = 'tag';
      buttonTag.innerHTML = `${lol.tags[tag]}`;
      divTags.appendChild(buttonTag);
    }
  }
}

/**
 * дописать
 */


/**
 * однажды дописать
 */
function fillUser() {
  const divCrad = document.getElementById('cardID');

  const divCardMain = document.createElement('div');
  divCardMain.id = 'cardMainID';
  divCardMain.className = 'card-main-profile';
  divCrad.appendChild(divCardMain);

  fillUserProfile();

  const divEdit = document.createElement('div');
  divEdit.id = 'editID';
  divEdit.className = 'actions-container-profile';
  divCrad.appendChild(divEdit);

  fillEdit();
}

function fillEdit() {

  
  const divEdit = document.getElementById('editID');

  const buttonLogout = document.createElement('button');
  buttonLogout.className = 'menu-icon';
  divEdit.appendChild(buttonLogout);
  const imgLogout = document.createElement('img');
  imgLogout.src = 'icons/button_previous_white.svg';
  imgLogout.style.width = '50px';
  imgLogout.style.height = '50px';
  imgLogout.className = 'profile-logout';

  
  buttonLogout.appendChild(imgLogout);



  const buttonEdit = document.createElement('button');
  buttonEdit.className = 'menu-icon';
  divEdit.appendChild(buttonEdit);
  const imgEdit = document.createElement('img');
  imgEdit.src = 'icons/button_edit_white.svg';
  imgEdit.style.width = '50px';
  imgEdit.style.height = '50px';
  imgEdit.alt = 'edit';
  imgEdit.className = 'profile-edit';
  buttonEdit.appendChild(imgEdit);
  
  
}

function fillCard() {
  const divCrad = document.getElementById('cardID');

  const divCardMain = document.createElement('div');
  divCardMain.id = 'cardMainID';
  divCardMain.className = 'card-main-expand';
  divCrad.appendChild(divCardMain);

  fillCardMain();

  const divshrink = document.createElement('div');
  divshrink.id = 'shrinkID';
  divshrink.className = 'forshrink';
  divCrad.appendChild(divshrink);

  fillshrink();
}

function fillshrink() {
  const divshrink = document.getElementById('shrinkID');

  const buttonshrink = document.createElement('button');
  buttonshrink.className = 'menu-icon';
  divshrink.appendChild(buttonshrink);
  const imgshrink = document.createElement('img');
  imgshrink.src = 'icons/button_shrink_white.svg';
  imgshrink.className = 'shrink-card';
  imgshrink.style.width = '50px';
  imgshrink.style.height = '50px';
  imgshrink.alt = 'shrink';
  buttonshrink.appendChild(imgshrink);
}

/**
 * @param {String} icon - путь до иконки для кнопки с действием
 * @param {String} action - класс действия
 * @return {HTMLButtonElement} - полученная кнопка
 */
function createActionElement(icon, action) {
  const actionElement = document.createElement('button');
  actionElement.className = 'menu-icon';
  const Icon = document.createElement('img');
  Icon.src = icon;
  Icon.width = 40;
  Icon.height = 40;
  Icon.classList.add(action);

  actionElement.appendChild(Icon);

  return actionElement;
}




/**
 * Отрисовывает меню с заданным активным айтемом
 * @param {String} activeItem - активный пункт меню
 */
function addMenu(activeItem) {
  const menu = createElementWithClass('div', 'tapbar-container');
  if (activeItem === 'feed') {
    menu.appendChild(createActionElement('icons/tapbar_feed_white_selected.svg', 'menu-feed'));
  } else {
    menu.appendChild(createActionElement('icons/tapbar_feed_white_deselected.svg', 'menu-feed'));
  }
  if (activeItem === 'likes') {
    menu.appendChild(createActionElement('icons/tapbar_likes_white_selected.svg', 'menu-likes'));
  } else {
    menu.appendChild(createActionElement('icons/tapbar_likes_white_deselected.svg', 'menu-likes'));
  }
  if (activeItem === 'chat') {
    menu.appendChild(createActionElement('icons/tapbar_chat_white_selected.svg', 'menu-chat'));
  } else {
    menu.appendChild(createActionElement('icons/tapbar_chat_white_deselected.svg', 'menu-chat'));
  }
  if (activeItem === 'profile') {
    menu.appendChild(createActionElement('icons/tapbar_user_white_selected.svg', 'menu-profile'));
  } else {
    menu.appendChild(createActionElement('icons/tapbar_user_white_deselected.svg', 'menu-profile'));
  }

  root.appendChild(menu);
}

/**
 * Очищает root и лишние EventListener-ы
 */



function createElementWithClass(tag, className) {
  const element = document.createElement(tag);
  element.className = className;
  return element;
}