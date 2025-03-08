// Генерація полів для електроприймачів
function generateDeviceInputs() {
    const count = parseInt(document.getElementById('deviceCount').value);
    const container = document.getElementById('deviceInputs');
    container.innerHTML = '';
    
    // Створення таблиці для введення даних
    const table = document.createElement('table');
    table.className = 'table table-bordered table-striped';
    
    // Заголовок таблиці
    const thead = document.createElement('thead');
    const headerRow = document.createElement('tr');
    
    const headers = [
        'Найменування ЕП',
        'ККД (ηн)',
        'Коеф. потужності (cos φ)',
        'Напруга (Uн, кВ)',
        'Кількість (n, шт)',
        'Потужність (Pн, кВт)',
        'Коеф. використання (КВ)',
        'Коеф. реактивної потужності (tgφ)'
    ];
    
    headers.forEach(header => {
        const th = document.createElement('th');
        th.textContent = header;
        headerRow.appendChild(th);
    });
    
    thead.appendChild(headerRow);
    table.appendChild(thead);
    
    // Тіло таблиці
    const tbody = document.createElement('tbody');
    
    for (let i = 0; i < count; i++) {
        const row = document.createElement('tr');
        
        // Назва ЕП
        const nameCell = document.createElement('td');
        const nameInput = document.createElement('input');
        nameInput.type = 'text';
        nameInput.name = `name_${i}`;
        nameInput.className = 'form-control';
        nameInput.value = `Електроприймач ${i+1}`;
        nameCell.appendChild(nameInput);
        row.appendChild(nameCell);
        
        // ККД
        const effCell = document.createElement('td');
        const effInput = document.createElement('input');
        effInput.type = 'number';
        effInput.name = `efficiency_${i}`;
        effInput.className = 'form-control';
        effInput.min = '0';
        effInput.max = '1';
        effInput.step = '0.01';
        effInput.value = '0.92';
        effCell.appendChild(effInput);
        row.appendChild(effCell);
        
        // Коефіцієнт потужності
        const pfCell = document.createElement('td');
        const pfInput = document.createElement('input');
        pfInput.type = 'number';
        pfInput.name = `powerFactor_${i}`;
        pfInput.className = 'form-control';
        pfInput.min = '0';
        pfInput.max = '1';
        pfInput.step = '0.01';
        pfInput.value = '0.9';
        pfCell.appendChild(pfInput);
        row.appendChild(pfCell);
        
        // Напруга
        const voltCell = document.createElement('td');
        const voltInput = document.createElement('input');
        voltInput.type = 'number';
        voltInput.name = `voltage_${i}`;
        voltInput.className = 'form-control';
        voltInput.min = '0';
        voltInput.step = '0.01';
        voltInput.value = '0.38';
        voltCell.appendChild(voltInput);
        row.appendChild(voltCell);
        
        // Кількість
        const quantCell = document.createElement('td');
        const quantInput = document.createElement('input');
        quantInput.type = 'number';
        quantInput.name = `quantity_${i}`;
        quantInput.className = 'form-control';
        quantInput.min = '1';
        quantInput.value = (i % 4) + 1;
        quantCell.appendChild(quantInput);
        row.appendChild(quantCell);
        
        // Потужність
        const powerCell = document.createElement('td');
        const powerInput = document.createElement('input');
        powerInput.type = 'number';
        powerInput.name = `power_${i}`;
        powerInput.className = 'form-control';
        powerInput.min = '0';
        powerInput.step = '0.1';
        powerInput.value = (i+1) * 10;
        powerCell.appendChild(powerInput);
        row.appendChild(powerCell);
        
        // Коефіцієнт використання
        const usageCell = document.createElement('td');
        const usageInput = document.createElement('input');
        usageInput.type = 'number';
        usageInput.name = `usageFactor_${i}`;
        usageInput.className = 'form-control';
        usageInput.min = '0';
        usageInput.max = '1';
        usageInput.step = '0.01';
        usageInput.value = (0.15 + i*0.05).toFixed(2);
        usageCell.appendChild(usageInput);
        row.appendChild(usageCell);
        
        // Коефіцієнт реактивної потужності
        const reactCell = document.createElement('td');
        const reactInput = document.createElement('input');
        reactInput.type = 'number';
        reactInput.name = `reactiveFactor_${i}`;
        reactInput.className = 'form-control';
        reactInput.min = '0';
        reactInput.step = '0.01';
        reactInput.value = i % 2 ? '1.0' : '1.33';
        reactCell.appendChild(reactInput);
        row.appendChild(reactCell);
        
        tbody.appendChild(row);
    }
    
    table.appendChild(tbody);
    container.appendChild(table);
}

// Генерація полів для крупних електроприймачів
function generateLargeDeviceInputs() {
    const count = parseInt(document.getElementById('largeDeviceCount').value);
    const container = document.getElementById('largeDeviceInputs');
    container.innerHTML = '';
    
    if (count <= 0) {
        return;
    }
    
    // Створення таблиці для введення даних
    const table = document.createElement('table');
    table.className = 'table table-bordered table-striped';
    
    // Заголовок таблиці
    const thead = document.createElement('thead');
    const headerRow = document.createElement('tr');
    
    const headers = [
        'Найменування ЕП',
        'ККД (ηн)',
        'Коеф. потужності (cos φ)',
        'Напруга (Uн, кВ)',
        'Кількість (n, шт)',
        'Потужність (Pн, кВт)',
        'Коеф. використання (КВ)',
        'Коеф. реактивної потужності (tgφ)'
    ];
    
    headers.forEach(header => {
        const th = document.createElement('th');
        th.textContent = header;
        headerRow.appendChild(th);
    });
    
    thead.appendChild(headerRow);
    table.appendChild(thead);
    
    // Тіло таблиці
    const tbody = document.createElement('tbody');
    
    const defaultNames = ['Зварювальний трансформатор', 'Сушильна шафа'];
    const defaultPowers = [100, 120];
    const defaultUsage = [0.2, 0.8];
    const defaultReactive = [3, 0];
    
    for (let i = 0; i < count; i++) {
        const row = document.createElement('tr');
        
        // Назва ЕП
        const nameCell = document.createElement('td');
        const nameInput = document.createElement('input');
        nameInput.type = 'text';
        nameInput.name = `large_name_${i}`;
        nameInput.className = 'form-control';
        nameInput.value = i < defaultNames.length ? defaultNames[i] : `Крупний ЕП ${i+1}`;
        nameCell.appendChild(nameInput);
        row.appendChild(nameCell);
        
        // ККД
        const effCell = document.createElement('td');
        const effInput = document.createElement('input');
        effInput.type = 'number';
        effInput.name = `large_efficiency_${i}`;
        effInput.className = 'form-control';
        effInput.min = '0';
        effInput.max = '1';
        effInput.step = '0.01';
        effInput.value = '0.92';
        effCell.appendChild(effInput);
        row.appendChild(effCell);
        
        // Коефіцієнт потужності
        const pfCell = document.createElement('td');
        const pfInput = document.createElement('input');
        pfInput.type = 'number';
        pfInput.name = `large_powerFactor_${i}`;
        pfInput.className = 'form-control';
        pfInput.min = '0';
        pfInput.max = '1';
        pfInput.step = '0.01';
        pfInput.value = '0.9';
        pfCell.appendChild(pfInput);
        row.appendChild(pfCell);
        
        // Напруга
        const voltCell = document.createElement('td');
        const voltInput = document.createElement('input');
        voltInput.type = 'number';
        voltInput.name = `large_voltage_${i}`;
        voltInput.className = 'form-control';
        voltInput.min = '0';
        voltInput.step = '0.01';
        voltInput.value = '0.38';
        voltCell.appendChild(voltInput);
        row.appendChild(voltCell);
        
        // Кількість
        const quantCell = document.createElement('td');
        const quantInput = document.createElement('input');
        quantInput.type = 'number';
        quantInput.name = `large_quantity_${i}`;
        quantInput.className = 'form-control';
        quantInput.min = '1';
        quantInput.value = '2';
        quantCell.appendChild(quantInput);
        row.appendChild(quantCell);
        
        // Потужність
        const powerCell = document.createElement('td');
        const powerInput = document.createElement('input');
        powerInput.type = 'number';
        powerInput.name = `large_power_${i}`;
        powerInput.className = 'form-control';
        powerInput.min = '0';
        powerInput.step = '0.1';
        powerInput.value = i < defaultPowers.length ? defaultPowers[i] : 100;
        powerCell.appendChild(powerInput);
        row.appendChild(powerCell);
        
        // Коефіцієнт використання
        const usageCell = document.createElement('td');
        const usageInput = document.createElement('input');
        usageInput.type = 'number';
        usageInput.name = `large_usageFactor_${i}`;
        usageInput.className = 'form-control';
        usageInput.min = '0';
        usageInput.max = '1';
        usageInput.step = '0.01';
        usageInput.value = i < defaultUsage.length ? defaultUsage[i] : 0.5;
        usageCell.appendChild(usageInput);
        row.appendChild(usageCell);
        
        // Коефіцієнт реактивної потужності
        const reactCell = document.createElement('td');
        const reactInput = document.createElement('input');
        reactInput.type = 'number';
        reactInput.name = `large_reactiveFactor_${i}`;
        reactInput.className = 'form-control';
        reactInput.min = '0';
        reactInput.step = '0.01';
        reactInput.value = i < defaultReactive.length ? defaultReactive[i] : 1.0;
        reactCell.appendChild(reactInput);
        row.appendChild(reactCell);
        
        tbody.appendChild(row);
    }
    
    table.appendChild(tbody);
    container.appendChild(table);
}

// Виконання розрахунків
function performCalculation() {
    const form = document.getElementById('electricalLoadForm');
    const formData = new FormData(form);
    
    fetch('/calculate', {
        method: 'POST',
        body: formData
    })
    .then(response => response.json())
    .then(data => {
        displayResults(data);
    })
    .catch(error => {
        console.error('Error:', error);
        alert('Помилка при виконанні розрахунків: ' + error.message);
    });
}

// Відображення результатів розрахунків
function displayResults(data) {
    const groupResultsTable = document.getElementById('groupResultsTable');
    const totalResultsTable = document.getElementById('totalResultsTable');
    
    // Очищення таблиць
    groupResultsTable.innerHTML = '';
    totalResultsTable.innerHTML = '';
    
    // Заповнення таблиці груп
    const groups = data.Groups.slice(0, 3); // Перші три групи (ШР1, ШР2, ШР3)
    
    // Рядки параметрів для груп
    const groupParams = [
        { name: 'Груповий коефіцієнт використання', key: 'UsageFactor', format: value => value.toFixed(4) },
        { name: 'Ефективна кількість ЕП', key: 'EffectiveDeviceCount', format: value => Math.floor(value) },
        { name: 'Розрахунковий коефіцієнт активної потужності', key: 'PowerFactor', format: value => value.toFixed(2) },
        { name: 'Розрахункове активне навантаження (кВт)', key: 'ActivePower', format: value => value.toFixed(2) },
        { name: 'Розрахункове реактивне навантаження (квар)', key: 'ReactivePower', format: value => value.toFixed(2) },
        { name: 'Повна потужність (кВА)', key: 'TotalPowerApparent', format: value => value.toFixed(2) },
        { name: 'Розрахунковий груповий струм (А)', key: 'GroupCurrent', format: value => value.toFixed(2) }
    ];
    
    groupParams.forEach(param => {
        const row = document.createElement('tr');
        
        const nameCell = document.createElement('td');
        nameCell.textContent = param.name;
        row.appendChild(nameCell);
        
        groups.forEach(group => {
            const valueCell = document.createElement('td');
            valueCell.textContent = param.format(group[param.key]);
            row.appendChild(valueCell);
        });
        
        groupResultsTable.appendChild(row);
    });
    
    // Заповнення таблиці загальних результатів
    const totalParams = [
        { name: 'Загальна кількість ЕП', key: 'TotalQuantity' },
        { name: 'Коефіцієнт використання цеху', key: 'UsageFactor', format: value => value.toFixed(2) },
        { name: 'Ефективна кількість ЕП цеху', key: 'EffectiveDeviceCount', format: value => Math.floor(value) },
        { name: 'Розрахунковий коефіцієнт активної потужності', key: 'PowerFactor', format: value => value.toFixed(2) },
        { name: 'Розрахункове активне навантаження (кВт)', key: 'ActivePower', format: value => value.toFixed(2) },
        { name: 'Розрахункове реактивне навантаження (квар)', key: 'ReactivePower', format: value => value.toFixed(2) },
        { name: 'Повна потужність (кВА)', key: 'TotalPowerApparent', format: value => value.toFixed(2) },
        { name: 'Розрахунковий груповий струм (А)', key: 'TotalCurrent', format: value => value.toFixed(2) }
    ];
    
    totalParams.forEach(param => {
        const row = document.createElement('tr');
        
        const nameCell = document.createElement('td');
        nameCell.textContent = param.name;
        row.appendChild(nameCell);
        
        const valueCell = document.createElement('td');
        const value = data[param.key];
        valueCell.textContent = param.format ? param.format(value) : value;
        row.appendChild(valueCell);
        
        totalResultsTable.appendChild(row);
    });
    
    // Показати контейнер результатів
    document.getElementById('resultContainer').style.display = 'block';
    
    // Прокрутка до результатів
    document.getElementById('resultContainer').scrollIntoView({ behavior: 'smooth' });
}

// Ініціалізація при завантаженні сторінки
document.addEventListener('DOMContentLoaded', function() {
    // Прив'язка обробників подій
    document.getElementById('generateDeviceInputs').addEventListener('click', generateDeviceInputs);
    document.getElementById('generateLargeDeviceInputs').addEventListener('click', generateLargeDeviceInputs);
    document.getElementById('calculateButton').addEventListener('click', performCalculation);
    
    // Автоматичне створення полів вводу при завантаженні
    generateDeviceInputs();
    generateLargeDeviceInputs();
});