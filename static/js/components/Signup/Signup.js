
export default class SignupComponent {
    #parent
    #data

    constructor(parent) {
        this.#parent = parent;
    }

    set data(data) {
        this.#data = data;
    }

    #renderDOM() {
        root.innerHTML = '';
    }

    render() {
        this.#renderDOM();
    }
}