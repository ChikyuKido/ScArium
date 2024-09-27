const openModalButton = document.getElementById('createAccount');
const closeModalButton = document.getElementById('closeModal');
const popupForm = document.getElementById('popupForm');
const typeSelect = document.getElementById('typeSelect');
const moodleFields = document.getElementById('moodleFields');
const digi4schoolFields = document.getElementById('digi4schoolFields');
const uploadButton = document.getElementById('uploadButton');
const fileInput = document.getElementById('displayImage');
const fileNameSpan = document.getElementById('fileName');
const buttonField = document.getElementById('uploadImageButtonField');

document.addEventListener("DOMContentLoaded", function () {
    initializePage();

    uploadButton.addEventListener('click', function() {
        fileInput.click();
    });
    fileInput.addEventListener('change', function() {
        const file = fileInput.files[0];
        if (file.size > 200 * 1024) {
            fileNameSpan.textContent = 'File is too large. Max size is 200KB.';
            fileInput.value = '';
        } else {
            fileNameSpan.textContent = `Selected file: ${file.name}`;
        }
    });
    openModalButton.addEventListener('click', () => {
        popupForm.style.display = 'block';
    });

    closeModalButton.addEventListener('click', () => {
        popupForm.style.display = 'none';
    });

    typeSelect.addEventListener('change', () => {
        const selectedType = typeSelect.querySelector('select').value;
        toggleFields(selectedType);
    });

    document.getElementById('submitBtn').addEventListener('click', () => {
        const selectedType = typeSelect.querySelector('select').value;
        handleAccountCreation(selectedType);
    });

    document.getElementById('cancelBtn').addEventListener('click', () => {
        popupForm.style.display = 'none';
    });
});

function initializePage() {
    displayMoodleAccounts();
    displayD4sAccounts();
}

function toggleFields(selectedType) {
    moodleFields.classList.add('is-hidden');
    digi4schoolFields.classList.add('is-hidden');
    if (selectedType === 'moodle') {
        moodleFields.classList.remove('is-hidden');
    } else if (selectedType === 'digi4school') {
        digi4schoolFields.classList.remove('is-hidden');
    }
    buttonField.classList.remove('is-hidden')
}

async function handleAccountCreation(selectedType) {
    const fields = selectedType === 'moodle' ? getMoodleFields() : getD4sFields();
    const errorMessage = document.getElementById('errorMessage');

    const error = await (selectedType === 'moodle' ? createMoodleAccount(...fields) : createD4sAccount(...fields));

    if (error) {
        errorMessage.textContent = error;
        errorMessage.classList.remove('is-hidden');
    } else {
        console.log("Account created successfully");
        document.getElementById('popupForm').style.display = 'none';
    }
}

function getMoodleFields() {
    return [
        document.getElementById('moodleUsername').value,
        document.getElementById('moodlePassword').value,
        document.getElementById('moodleDisplayName').value,
        document.getElementById('moodleInstanceUrl').value
    ];
}

function getD4sFields() {
    return [
        document.getElementById('digi4schoolUsername').value,
        document.getElementById('digi4schoolPassword').value,
        document.getElementById('digi4schoolDisplayName').value
    ];
}



async function createAccount(endpoint, accountData) {
    const { username, password, display_name, instance_url,image_file} = accountData;

    if (!username || !password || !display_name || (instance_url !== undefined && !instance_url)) {
        return "All fields must be filled";
    }

    if(image_file) {
        accountData.display_image = await convertToBase64(image_file);
    }else {
        accountData.display_image = "empty";
    }
    delete accountData.image_file;

    let errorString = "";
    const response = await fetch(endpoint, {
        method: "POST",
        body: JSON.stringify(accountData)
    });

    if (response.status !== 201) {
        const data = await response.json();
        errorString = data.error || 'Failed to create account';
    }
    //clear filename
    fileNameSpan.textContent = '';
    fileInput.value = '';
    return errorString;
}

function convertToBase64(file) {
    return new Promise((resolve, reject) => {
        const reader = new FileReader();
        reader.readAsDataURL(file);

        reader.onload = () => resolve(reader.result.split(',')[1]); // Get Base64 part only
        reader.onerror = (error) => reject(error);
    });
}

async function createMoodleAccount(username, password, displayName, instanceUrl) {
    return createAccount("/api/v1/account/createMoodleAccount", { username, password, display_name: displayName,
        image_file: fileInput.files.length === 0 ? null : fileInput.files[0], instance_url: instanceUrl });
}

async function createD4sAccount(username, password, displayName) {
    return createAccount("/api/v1/account/createD4sAccount", { username, password, display_name: displayName,
        image_file: fileInput.files.length === 0 ? null : fileInput.files[0]});
}

function displayMoodleAccounts() {
    fetchAccounts("/api/v1/account/getMoodleAccounts", createMoodleCard, 'moodleAccountsContainer');
}

function displayD4sAccounts() {
    fetchAccounts("/api/v1/account/getD4sAccounts", createD4sCard, 'd4sAccountsContainer');
}

function fetchAccounts(apiEndpoint, cardCreator, containerId) {
    const container = document.getElementById(containerId);

    fetch(apiEndpoint)
        .then(response => {
            if (response.status !== 200) {
                return response.json().then(data => {
                    throw new Error(data.error || 'Failed to retrieve accounts');
                });
            }
            return response.json();
        })
        .then(data => {
            data.forEach(account => {
                if(account.image_id !== "") {
                    account.image_url = "/api/v1/account/accountImage/"+account.image_id;
                }
                delete account.image_id;
                const card = cardCreator(account);
                container.appendChild(card);
            });
        })
        .catch(reason => {
            console.log(reason);
        });
}

function createMoodleCard({ image_url, display_name, username, instance_url, id }) {
    return createCard(image_url || "/imgs/moodle.png", display_name, username, instance_url, id);
}

function createD4sCard({ image_url, display_name, username, id }) {
    return createCard(image_url || "/imgs/d4s.png", display_name, username, null, id);
}

function createCard(imageSrc, displayName, username, instanceUrl, id) {
    const domain = instanceUrl ? getDomain(instanceUrl) : '';
    const cardHTML = `
        <div class="column is-one-quarter">
            <a href="/account/moodle/${id}" class="card-link">
                <div class="card" style="max-width: 300px;">
                    <div class="card-image">
                        <figure class="image is-4by3" style="margin: 0;">
                            <img src="${imageSrc}" alt="User Image" class="is-fullwidth" style="object-fit: fill;">
                        </figure>
                    </div>
                    <div class="card-content">
                        <div class="media">
                            <div class="media-content">
                                <p class="title is-6">${displayName}</p>
                                <p class="subtitle is-7">${username}</p>
                                ${instanceUrl ? `<p><a href="${instanceUrl}">${domain}</a></p>` : ''}
                            </div>
                        </div>
                    </div>
                </div>
            </a>
        </div>
    `;
    const tempDiv = document.createElement('div');
    tempDiv.innerHTML = cardHTML;
    return tempDiv.children[0];
}

function getDomain(url) {
    return new URL(url).hostname;
}
