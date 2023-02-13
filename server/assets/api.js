function updateTableList() {
    fetch('/api/metrics')
      .then(response => response.json())
      .then(retData => {
        const tableList = document.querySelector('#tables');
        tableList.innerHTML = '';
        console.log(tables);
        for (const ep of retData.buckets) {

          const row = document.createElement('tr');
          
          // Pridani sloupce s URL
          const nameCell = document.createElement('td');
          nameCell.innerText = ep.Url;
          row.appendChild(nameCell);

          // Pridani sloupce se statusem
          const statusCell = document.createElement('td');
            const statusBar = document.createElement('div');
            statusBar.id = "status_bar";

            // definice kde se zacal endpoint merit
            let start = Math.floor(Math.random() * 10) + 3;

            console.log(ep);
            console.log(retData);

            for (const dat of retData.graphData[ep.Url].slice(-30)) {
              const barItem = document.createElement('div');

              const sc = dat.StatusCode;
              console.log(sc);
              switch (true) {
                case (sc < 200):
                  barItem.className = "bar no-data";
                  break;
                case (199 < sc && sc < 400):
                  barItem.className = "bar ok";
                  break;
                case (399 < sc && sc < 500):
                  barItem.className = "bar warning";
                  break;
                case (sc > 499):
                  barItem.className = "bar critical";
                  break;
                default:
                  barItem.className = "bar";
              }
              barItem.title = "Status: " + dat.StatusCode + ", Time: " + dat.Time + ")";
              
              statusBar.appendChild(barItem);
            }
            

            statusCell.appendChild(statusBar);


          row.appendChild(statusCell);

          // Pridani sloupce s Akcemi
          const deleteCell = document.createElement('td');
            const deleteButton = document.createElement('button')
                deleteButton.className = "btn-danger btn";
                deleteButton.innerText = 'x';
                deleteButton.addEventListener('click', () => {
                    fetch(`/api/metrics/${ep.name}`, { method: 'DELETE' })
                        .then(() => updateTableList());
                    });
          
          deleteCell.appendChild(deleteButton);
          row.appendChild(deleteCell);

          // Cely radek pridame do tabulky
          tableList.appendChild(row);
          
        }
      });
  }

  updateTableList();

  const addTableForm = document.querySelector('#add-table-form');
  addTableForm.addEventListener('submit', event => {
    event.preventDefault();
    const tableName = document.querySelector('#table-name').value;
    const retentionDays = 15;

    let data = { "url": tableName};
    console.log("POST to /api/metrics data:", data);
    
    fetch('/api/metrics', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(data)
    })
    .then(() => updateTableList());
  });
