function createInputEdit(type, name, place, className) {
    const input = document.createElement('input');
    input.type = type;
    input.name = name;
    input.placeholder = place;
    input.className = className;

    return input
}

const fillForm = form => {
    const divName = document.createElement('div');
    divName.className = 'inputEdit'
    const inputName = createInputEdit('text', 'name', 'Elon', 'form-field text-without-icon');
    divName.appendChild(inputName);

    const divDate = document.createElement('div');
    divDate.className = 'inputEdit'
    const inputDate = createInputEdit('date', 'calendar', '20.06.2001', 'form-field text-with-icon');
    divDate.appendChild(inputDate);

    const divDesc = document.createElement('div');
    divDesc.className = 'inputEdit'
    const desc = document.createElement('textarea');
    // desc.type = 'text';
    desc.className = 'form-field-desc text-desc'
    desc.textContent = `Инст: @elonmask
Меня зовут Ильягу всем привет, я люблю аниме и особенно ван пис. Люблю Дензела Карри и Френка Оушена.
`
    divDesc.appendChild(desc);

    const divTags = document.createElement('div');
    divTags.className = 'inputEdit';

    const TagsContainer = document.createElement('div');
    TagsContainer.className = 'tag-container-edit';
    TagsContainer.id = 'tagsID';

    const buttonAddTags = document.createElement('button');
    buttonAddTags.id = 'addID';
    buttonAddTags.className = 'add';

    divTags.appendChild(TagsContainer);
    divTags.appendChild(buttonAddTags);

    const divSelect = document.createElement('div');
    divSelect.className = 'selectBox';

    const divImgs = document.createElement('div');
    divImgs.className = 'inputEdit';
    divImgs.innerHTML = `
<div class="im-container">
    <div style="position: relative;">
        <img src="../static/img/Elon_Musk_2015.jpg" class="im">
        <button class=removeImg></button>
    </div>
    <div style="position: relative;">
        <img src="../static/img/Elon_Musk_2015.jpg" class="im">
        <button class=removeImg></button>
    </div>
    <div style="position: relative;">
        <img src="../static/img/Elon_Musk_2015.jpg" class="im">
        <button class=removeImg></button>
    </div>
</div>
<button class="add"></button>
`
    form.appendChild(divName);
    form.appendChild(divDate);
    form.appendChild(divDesc);
    form.appendChild(divTags);
    form.appendChild(divSelect);
    form.appendChild(divImgs);
}

const fillSave = saveDiv => {
    const buttonSave = document.createElement('button');
    buttonSave.type = 'submit';
    buttonSave.className = 'login-button';

    const div = document.createElement('div');
    div.className = 'center-container';

    const span = document.createElement('span');
    span.className = 'login-button-text-edit';
    span.textContent = 'Сохранить';

    const imgNext = document.createElement('img');
    imgNext.src = '../static/icons/button_next_black.svg';
    imgNext.className = 'svg-next-edit'

    div.appendChild(span);
    div.appendChild(imgNext);
    buttonSave.appendChild(div);

    saveDiv.appendChild(buttonSave);
}

function renderEdit() {
    const root = document.getElementById('root');

    const form = document.createElement('form');
    form.className = 'edit-form';
    const saveDiv = document.createElement('div');
    saveDiv.className = 'center-container';
    saveDiv.style = 'padding-top: 20%;';
    root.appendChild(form);
    root.appendChild(saveDiv);

    fillForm(form);
    fillSave(saveDiv);

    const divSelectBox = document.getElementsByClassName('selectBox')[0];

    const selectBoxItems = ['anime', 'gaming', 'soccer', 'music'];
    const existsSelectBoxItems = []

    const selectBox = document.createElement('select');
    const selectItem = document.createElement('option');
    selectItem.textContent = 'Тэги';
    selectItem.value = 'Тэги';
    selectItem.disabled = true;
    selectBox.appendChild(selectItem);
    selectBoxItems.forEach(function (item, i, selectBoxItems) {
        const selectItem = document.createElement('option');
        selectItem.textContent = item;
        selectItem.value = item;
        selectBox.appendChild(selectItem);
    });

    const tagsCont = document.getElementById('tagsID');

    selectBox.onchange = function () {
        divSelectBox.innerHTML = '';
        const {
            value
        } = selectBox;
        if (existsSelectBoxItems.indexOf(value) != -1) {
            return
        }
        const tag = document.createElement('div');
        tag.className = 'tag-edit';
        tag.textContent = value;
        tagsCont.appendChild(tag);
        tag.disabled = true;
        existsSelectBoxItems.push(value)
    }

    root.addEventListener('click', function (e) {
        e.preventDefault();
        const {
            target
        } = e;

        if (target.tagName.toLowerCase() === 'option') {
            divSelectBox.innerHTML = '';
        }
        if (target.id === 'addID') {
            divSelectBox.appendChild(selectBox);
        }
    });
}

renderEdit();