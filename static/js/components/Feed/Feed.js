


let currentCard;
let previousCard;
let previousCard2;


let x
let y
let x1
let y1


function handleTouchEnd(event) {
    if (!x1 || !y1) {
      return;
    }
  
    if (x === null) {
      x = x1;
    }
    if (x1 - x < -200) {
      currentCard.style.animation = 'liked 1s ease 1';
      const cardToRemove = currentCard
      setTimeout(remove(cardToRemove), 1000);
      
      
      const {id} = window.Feed.getCurrentProfile()
      console.log(id)
      setTimeout(window.Feed.feedGet(id), 1000);
      
        const feed = new FeedComponent(root);
        feed.render();
  
      x1 = null;
      x = null;
    } else if (x1 - x > 200) {
      currentCard.style.animation = 'liked 1s ease 1';
      const cardToRemove = currentCard
      setTimeout(remove(cardToRemove), 1000);
      const {id} = window.Feed.getCurrentProfile()
   
      console.log(id)
      setTimeout(window.Feed.feedGet(id), 1000);
      
        const feed = new FeedComponent(root);
        feed.render();
      
      x1 = null;
      x = null;
    } else {
      const {
        target
      } = event;
      if (!(target.class === 'expand-class' || target.alt === 'shrink')) {
        previousCard.style.animation = 'shrinkSecondary 1s linear 1';
        previousCard2.style.animation = 'shrinkThird 1s linear 1';
        currentCard.style.animation = 'spin2 1s linear 1';
        setTimeout(returnToStart, 1000);
      }
    }
  }


function remove(cardToRemove) {
    cardToRemove.style.opacity = 0;
  }


function handleTouchStart(event) {
    const {
      touches
    } = event;
    currentCard.style.animation = '';
    x1 = touches[0].clientX;
    y1 = touches[0].clientY;
  }

  /**
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
  
function handleTouchMove(event) {
    const {
      touches
    } = event;
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
   */


export default class FeedComponent {
    #parent
    #data


    constructor(parent) {
        this.#parent = parent;
    }

    set data(data) {
        this.#data = data;
    }
    
    /**
 * @param {String} tag - название тэга
 * @param {String} className - название класса
 * @return {HTMLElement} - Полученный элемент с тэгом
 */
   #createElementWithClass(tag, className) {
    const element = document.createElement(tag);
    element.className = className;
    return element;
  }


    #renderDOM() {
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
        root.innerHTML = '';
        
        document.addEventListener('touchstart', handleTouchStart, false);
        document.addEventListener('touchmove', handleTouchMove, false);
        document.addEventListener('touchend', handleTouchEnd, false);
        
        const currentobj = window.Feed.getCurrentProfile();
        console.log(currentobj);
        if(!currentobj) {
          root.innerHTML='';
          const outOfCards = this.#createElementWithClass('div','out-of-cards');
          outOfCards.innerText = 'Карточки кончились'
          root.appendChild(outOfCards);
          
          return
        } 
        const card = this.#createElementWithClass('div', 'card-main');
        const image = document.createElement('img');
        image.src = currentobj.photoSrc;
        image.className = 'profile-image';
        card.appendChild(image);
        const bottomPanel = this.#createElementWithClass('div', 'bottom-panel');
        const nameContainer = this.#createElementWithClass('div', 'name-container');
        const name = this.#createElementWithClass('div', 'name');
        name.innerText = currentobj.firstName;
        const age = this.#createElementWithClass('div', 'age');
        age.innerText = currentobj.age;
        nameContainer.appendChild(name);
        nameContainer.appendChild(age);
        bottomPanel.appendChild(nameContainer);
        const actionsContainer = this.#createElementWithClass('div', 'actions-container');
      
      
        actionsContainer.appendChild(createActionElement('icons/button_dislike_white.svg', 'dislike-card'));
      
        actionsContainer.appendChild(createActionElement('icons/button_expand_white.svg', 'expand-card'));
      
        actionsContainer.appendChild(createActionElement('icons/tapbar_likes_white_selected.svg', 'like-card'));
      
      
        bottomPanel.appendChild(actionsContainer);
        card.appendChild(bottomPanel);
      
      
        root.appendChild(this.#createElementWithClass('div', 'card3'));
        root.appendChild(this.#createElementWithClass('div', 'card3'));
        root.appendChild(this.#createElementWithClass('div', 'card2'));
        const cardNew = this.#createElementWithClass('div', 'card');
        cardNew.appendChild(card);
        root.appendChild(cardNew);
      
        currentCard = document.getElementsByClassName('card')[0];
        console.log(currentCard);
        previousCard = document.getElementsByClassName('card2')[0];
        previousCard2 = document.getElementsByClassName('card3')[1];
    }
    render() {
        this.#renderDOM();
    }
   
     
    
      
      
}