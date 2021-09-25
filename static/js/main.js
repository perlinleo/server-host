const root = document.getElementById('root');

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
  createProfile: {
    href: '/createProfile',
    name: 'Создать профиль',
    open: createProfilePage,
  },
  profile: {
    href: '/profile',
    name: 'Профиль',
    open: userProfileRender,
  },
  'menu-feed': {
    link: '/feed',
    name: 'feed',
    open: renderFeed,
  },
  'menu-profile': {
    link: '/profile',
    name: 'profile',
    open: userProfileRender,
  },
  'expand-card': {
    link: '/feed',
    name: 'feed',
    open: profileRender,
  },
  'shrink-card': {
    link: '/feed',
    name: 'feed',
    open: renderFeed,
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
    open: logoutCookie,
  }
}

const emailRegExp = /^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9-]+(?:\.[a-zA-Z0-9-]+)*$/;
const passwordRegExp = /^(?=.*\d)(?=.*[a-z])(?=.*[A-Z])(?=.*[a-zA-Z]).{8,}$/;


function createInput(type, text, name) {
  const input = document.createElement('input');
  input.type = type;
  input.name = name;
  input.placeholder = text;
  // input.classList.add('form-field');

  return input;
}

function createCenterContainer() {
  const divContainer = document.createElement('div');
  divContainer.classList.add('center-container');

  return divContainer;
}


function loginPageError(error){
  errorField = document.getElementsByClassName('login-error')[0]
  errorField.style.visibility="visible"
  errorField.innerHTML=error
}

function loginPage() {
  root.innerHTML = '';

  // --------------------------------------------------------
  window.addEventListener('load', (e) => {
      e.preventDefault();
      loginWithCookie()
  })
    // --------------------------------------------------------

  const header = createCenterContainer();

  const headerText = document.createElement('span');
  headerText.textContent = 'Войти';
  headerText.classList.add('login-header');

  header.appendChild(headerText);
  root.appendChild(header);

  const form = document.createElement('form');
  form.classList.add('login-form');
  const errorField = createElementWithClass('div', 'login-error');
  errorField.innerHTML="error placeholder";
  const emailInput = createInput('email', 'Почта', 'email');
  emailInput.addEventListener('input', () => {
    const test = emailInput.value.length === 0 || emailRegExp.test(emailInput.value);
    if (test) {
      emailInput.className = 'form-field-valid';
    } else {
      emailInput.className = 'form-field-novalid'
    }
  })

  const passwordInput = createInput('password', 'Пароль', 'password');
  passwordInput.addEventListener('input', () => {
    const test = passwordInput.value.length === 0 || passwordRegExp.test(passwordInput.value);

    if (test) {
      passwordInput.className = 'form-field-valid';
    } else {
      passwordInput.className = 'form-field-novalid'
    }
  })

  window.addEventListener('load', () => {
    const testEmail = emailInput.value.length === 0 || emailRegExp.test(emailInput.value);
    emailInput.className = testEmail ? 'form-field-valid' : 'form-field-novalid';
    const testPassword = passwordInput.value.length === 0 || passwordRegExp.test(passwordInput.value);
    passwordInput.className = testPassword ? 'form-field-valid' : 'form-field-novalid';
  })

  // кнопка войти
  const submitButton = document.createElement('button');
  submitButton.type = 'submit';
  submitButton.classList.add('login-button');

  const buttonFilling = createCenterContainer();
  const buttonText = document.createElement('span');
  buttonText.textContent = 'Войти';
  buttonText.classList.add('login-button-text');
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

  const emailFieldWithIcon = document.createElement('div');
  emailFieldWithIcon.classList.add('input-with-icon');
  const passwordFieldWithIcon = document.createElement('div');
  passwordFieldWithIcon.classList.add('input-with-icon');

  const logoBg = document.createElement('div');
  logoBg.classList.add('drip-logo-bg');

  const formContainer = createCenterContainer();


  emailFieldWithIcon.appendChild(emailInput);
  emailFieldWithIcon.appendChild(emailIcon);
  passwordFieldWithIcon.appendChild(passwordInput);
  passwordFieldWithIcon.appendChild(passwordIcon);

  logoBg.appendChild(errorField);
  logoBg.appendChild(emailFieldWithIcon);
  logoBg.appendChild(passwordFieldWithIcon);

  form.appendChild(logoBg);
  form.appendChild(submitButton);

  formContainer.appendChild(form);

  const regLinkContainer = createCenterContainer();
  const regLink = document.createElement('a');
  regLink.classList.add('reg-link');
  regLink.href = '/signup';
  regLink.textContent = 'Зарегестрироваться';
  regLink.dataset.section = 'signup';
  regLinkContainer.appendChild(regLink);

  form.addEventListener('submit', (e) => {
    e.preventDefault();
    const testEmail = emailRegExp.test(emailInput.value);
    const testPassword = passwordRegExp.test(passwordInput.value);

    if (!testEmail || !testPassword) {
      emailInput.className = 'form-field-novalid';
      passwordInput.className = 'form-field-novalid';
      // e.preventDefault();
      return;
    }

    const email = emailInput.value.trim();
    const password = passwordInput.value.trim();
    loginWithCredentials(email,password);
  })

  root.appendChild(formContainer);
  root.appendChild(regLinkContainer);
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
    
    const email = emailInput.value.trim();
    const password = passwordInput.value.trim();
 
    signupUser(email,password);
  })

  root.appendChild(formContainer);
}

function createProfilePage() {
  root.innerHTML = '<span>Create Profile Page</span>';
}



loginPage();


root.addEventListener('click', (e) => {
  const { target } = e;
  console.log(target.className);
  if (configApp[target.className]) {
    clearRoot();
    configApp[target.className].open();
    addMenu(configApp[target.className].name);
  }
  if (target instanceof HTMLAnchorElement) {
    e.preventDefault();

    configApp[target.dataset.section].open();
  }
})


// AHAHHAHAHAHAH


/**
    Обрабатывает нажатия
 * @param {event} event - Событие
 */
function clickButtons(event) {
  const { target } = event;
  if (configApp[target.className]) {
    clearRoot();
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

let currentCard;
let previousCard;
let previousCard2;


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
function userProfileRender() {
  const divCrad = document.createElement('div');
  divCrad.id = 'cardID';
  divCrad.className = 'card-expand';
  root.appendChild(divCrad);

  const divTapBar = document.createElement('div');
  divTapBar.id = 'tapbarID';
  divTapBar.className = 'tapbar-container';
  root.appendChild(divTapBar);

  fillUser();
}

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

/**
//  __  __                   _  __
// |  \/  |   ___    _ __   | |/ /   ___   _   _   ___
// | |\/| |  / _ \  | '_ \  | ' /   / _ \ | | | | / __|
// | |  | | | (_) | | | | | | . \  |  __/ | |_| | \__ \
// |_|  |_|  \___/  |_| |_| |_|\_\  \___|  \__, | |___/
//                                         |___/
 */
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
  imgLogout.alt = 'logout';
  
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

/**
//  __  __                   _  __
// |  \/  |   ___    _ __   | |/ /   ___   _   _   ___
// | |\/| |  / _ \  | '_ \  | ' /   / _ \ | | | | / __|
// | |  | | | (_) | | | | | | . \  |  __/ | |_| | \__ \
// |_|  |_|  \___/  |_| |_| |_|\_\  \___|  \__, | |___/
//                                         |___/
 */
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

/**
//  __  __                   _  __
// |  \/  |   ___    _ __   | |/ /   ___   _   _   ___
// | |\/| |  / _ \  | '_ \  | ' /   / _ \ | | | | / __|
// | |  | | | (_) | | | | | | . \  |  __/ | |_| | \__ \
// |_|  |_|  \___/  |_| |_| |_|\_\  \___|  \__, | |___/
//                                         |___/
 */
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
//  __  __                   _  __
// |  \/  |   ___    _ __   | |/ /   ___   _   _   ___
// | |\/| |  / _ \  | '_ \  | ' /   / _ \ | | | | / __|
// | |  | | | (_) | | | | | | . \  |  __/ | |_| | \__ \
// |_|  |_|  \___/  |_| |_| |_|\_\  \___|  \__, | |___/
//                                         |___/
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
 *
 * @param {String} tag - название тэга
 * @param {String} className - название класса
 * @return {HTMLElement} - Полученный элемент с тэгом
 */
function createElementWithClass(tag, className) {
  const element = document.createElement(tag);
  element.className = className;
  return element;
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
 *
 * Очищает root и лишние EventListener-ы
 */
function clearRoot() {
  root.innerHTML = '';
  document.removeEventListener('touchstart', handleTouchStart);
  document.removeEventListener('touchmove', handleTouchMove);
  document.removeEventListener('touchend', handleTouchEnd);
}




/**
 * Рендерит ленту с профилями
 */
function renderFeed() {
  document.addEventListener('click', clickButtons, false);
  document.addEventListener('touchstart', handleTouchStart, false);
  document.addEventListener('touchmove', handleTouchMove, false);
  document.addEventListener('touchend', handleTouchEnd, false);
  const currentobj = sample[counter];
  if(!currentobj) {
    root.innerHTML='';
    const outOfCards = createElementWithClass('div','out-of-cards');
    outOfCards.innerText = 'Карточки кончились'
    root.appendChild(outOfCards);
    addMenu('feed');
    return
  } 
  const card = createElementWithClass('div', 'card-main');
  const image = document.createElement('img');
  image.src = currentobj.photoSrc;
  image.className = 'profile-image';
  card.appendChild(image);
  const bottomPanel = createElementWithClass('div', 'bottom-panel');
  const nameContainer = createElementWithClass('div', 'name-container');
  const name = createElementWithClass('div', 'name');
  name.innerText = currentobj.firstName;
  const age = createElementWithClass('div', 'age');
  age.innerText = currentobj.age;
  nameContainer.appendChild(name);
  nameContainer.appendChild(age);
  bottomPanel.appendChild(nameContainer);
  const actionsContainer = createElementWithClass('div', 'actions-container');


  actionsContainer.appendChild(createActionElement('icons/button_dislike_white.svg', 'dislike-card'));

  actionsContainer.appendChild(createActionElement('icons/button_expand_white.svg', 'expand-card'));

  actionsContainer.appendChild(createActionElement('icons/tapbar_likes_white_selected.svg', 'like-card'));


  bottomPanel.appendChild(actionsContainer);
  card.appendChild(bottomPanel);


  root.appendChild(createElementWithClass('div', 'card3'));
  root.appendChild(createElementWithClass('div', 'card3'));
  root.appendChild(createElementWithClass('div', 'card2'));
  const cardNew = createElementWithClass('div', 'card');
  cardNew.appendChild(card);
  root.appendChild(cardNew);

  currentCard = document.getElementsByClassName('card')[0];
  previousCard = document.getElementsByClassName('card2')[0];
  previousCard2 = document.getElementsByClassName('card3')[1];
}

/**
 * Рендерит следующую карточку
 */
function nextCharacter() {
  const currentobj = sample[counter];
  if(!currentobj){
    const card1 = document.getElementsByClassName('card3')[0]
    if (card1) card1.style.animation='liked 1s ease 1'
    const card2 = document.getElementsByClassName('card3')[1]
    if (card2) card2.style.animation='liked 1s ease 1'
    const card3 = document.getElementsByClassName('card2')[0]
    if (card3) card3.style.animation='liked 1s ease 1'
    const card4 = document.getElementsByClassName('card')[0]
    if(card4) card4.style.animation='liked 1s ease 1'
    setTimeout(()=> {
      root.innerHTML='';
      const outOfCards = createElementWithClass('div','out-of-cards');
      outOfCards.innerText = 'Карточки кончились'
      root.appendChild(outOfCards);
      addMenu('feed');
      
    },1000);
    return;
  }
  const card = createElementWithClass('div', 'card-main');
  const image = document.createElement('img');
  image.src = currentobj.photoSrc;
  image.className = 'profile-image';
  card.appendChild(image);
  const bottomPanel = createElementWithClass('div', 'bottom-panel');
  const nameContainer = createElementWithClass('div', 'name-container');
  const name = createElementWithClass('div', 'name');
  name.innerText = currentobj.firstName;
  const age = createElementWithClass('div', 'age');
  age.innerText = currentobj.age;
  nameContainer.appendChild(name);
  nameContainer.appendChild(age);
  bottomPanel.appendChild(nameContainer);
  const actionsContainer = createElementWithClass('div', 'actions-container');


  actionsContainer.appendChild(createActionElement('icons/button_dislike_white.svg', 'dislike-card'));

  actionsContainer.appendChild(createActionElement('icons/button_expand_white.svg', 'expand-card'));

  actionsContainer.appendChild(createActionElement('icons/tapbar_likes_white_selected.svg', 'like-card'));


  bottomPanel.appendChild(actionsContainer);
  card.appendChild(bottomPanel);
  root.appendChild(card);
  
  const mainCard = document.getElementsByClassName('card2')[0];
  if (mainCard) {
    mainCard.className = 'card';
    mainCard.appendChild(card);
  } else {
    const mainCard = document.createElement('div');
    mainCard.className = 'card';
    mainCard.appendChild(card);
  }

  const card3old = document.getElementsByClassName('card3')[1];
  if (card3old) {
    card3old.className = 'card2';
  }
  const card1 = document.getElementsByClassName('card')[0];
  if (card1) {
    root.removeChild(document.getElementsByClassName('card')[1]);
  }

  root.innerHTML = '<div class="card3"></div>' + root.innerHTML;

  
  currentCard = document.getElementsByClassName('card')[0];
  previousCard = document.getElementsByClassName('card2')[0];
  previousCard2 = document.getElementsByClassName('card3')[1];
  const cardMain = document.getElementsByClassName('card-main')[0];
  if (cardMain) {
    cardMain.style.animation = 'appearance 0.3s linear 1';
  }
}




let x1 = null;
let y1 = null;

let x = null;
let y = null;

/**
 * Убирает текущую карточку
 */
function remove() {
  currentCard.style.opacity = 0;
}

/**
 *
 * @param {Event} event - начало нажатия
 */
function handleTouchStart(event) {
  const { touches } = event;
  currentCard.style.animation = '';
  x1 = touches[0].clientX;
  y1 = touches[0].clientY;
}

/**
 *
 * @param {HTMLElement} element - элемент для движения
 * @param {*} diffX - движение по оси X
 * @param {*} diffY - движение по оси Y
 */
function moveRight(element, diffX, diffY) {
  element.style.transform = `translate(${diffX}px, ${diffY}px)`;
  const topScale = 12 - Math.abs(diffX / 40);
  if (topScale > 5) {
    previousCard.style.top = `${topScale}%`;
  }
  const heightScale = 75 + Math.abs(diffX / 40);
  if (heightScale < 80) {
    previousCard.style.height = `${heightScale}%`;
  }

  const widthScale = 90 + Math.abs(diffX / 40);
  if (widthScale < 95) {
    previousCard.style.width = `${widthScale}%`;
  }

  const topScale2 = 19 - Math.abs(diffX / 40);
  if (topScale2 > 12) {
    previousCard2.style.top = `${topScale2}%`;
  }
  const heightScale2 = 70 + Math.abs(diffX / 40);
  if (heightScale2 < 75) {
    previousCard2.style.height = `${heightScale2}%`;
  }

  const widthScale2 = 85 + Math.abs(diffX / 40);
  if (widthScale2 < 90) {
    previousCard2.style.width = `${widthScale2}%`;
  }


  element.style.transform += `rotateZ(${(diffX/10)}deg)`;
}

/**
 * Обработчик начала свайпа по текущей карточке
 * @param {Event} event - событие
 */
function handleTouchMove(event) {
  const { touches } = event;
  x = touches[0].clientX;
  y = touches[0].clientY;
  if (window.innerHeight < y || window.innerWidth < x || y < 0 || x < 0) {
    return;
  }
  moveRight(currentCard, x - x1, y - y1);
}

/**
 * Сбрасывает стили на всех карточках
 */
function returnToStart() {
  currentCard.style.transform = 'translate(0px, 0px)';
  previousCard.style.width = '90%';
  previousCard.style.height = '75%';
  previousCard.style.top = '12%';
  previousCard.style.animation = '';
  previousCard2.style.width = '85%';
  previousCard2.style.height = '70%';
  previousCard2.style.top = '19%';
  previousCard2.style.animation = '';
}


/**
 * Обработчик конца свайпа - здесь фиксируется,
 * Была ли карточка лайкнута, дизлайкнута или
 * ей просто повозили по экрану (по приколу)
 * @param {event} event - событие
 *
 */
function handleTouchEnd(event) {
  if (!x1 || !y1) {
    return;
  }

  if (x === null) {
    x = x1;
  }
  if (x1 - x < -200) {
    currentCard.style.animation = 'liked 1s ease 1';
    setTimeout(remove, 1000);
    setTimeout(nextCharacter, 1000);
    counter++;

    x1 = null;
    x = null;
  } else if (x1 - x > 200) {
    currentCard.style.animation = 'liked 1s ease 1';

    setTimeout(remove, 1000);
    setTimeout(nextCharacter, 1000);
    counter++;
    x1 = null;
    x = null;
  } else {
    const { target } = event;
    if (!(target.class === 'expand-class' || target.alt === 'shrink')) {
      previousCard.style.animation = 'shrinkSecondary 1s linear 1';
      previousCard2.style.animation = 'shrinkThird 1s linear 1';
      currentCard.style.animation = 'spin2 1s linear 1';
      setTimeout(returnToStart, 1000);
    }
  }
}
