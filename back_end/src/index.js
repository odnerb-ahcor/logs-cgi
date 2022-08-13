const serve = require('express')
const app = serve()
const port = 3000
const cors = require('cors')
const logs = require('./ler_logs')
const db = require('./logs')

app.use(cors())
app.use(serve.json());

app.get('/ler/:id?/:horas?', async (req, res) => {
  await logs.buscar_logs()

  const { id, horas } = req.params
  let vaidacao = false

  if (!(!db.data[id])){
    vaidacao = !(db.data[id].horas === horas)
  }

  if (!id || vaidacao)
    res.json(db.data)
  else
    res.json(logs.retornaLogs(id))  
})

app.get('/apagar', (req, res) => {
  db.data.length = 0
  db.data.id = 0
  res.send('Apagado')
})

app.listen(port, () => {
  console.log(`Example app listening on port ${port}`)
})