const serve = require('express')
const app = serve()
const port = 3000
const cors = require('cors')
const logs = require('./ler_logs')
const db = require('./logs')

app.use(cors())
app.use(serve.json());

app.get('/ler', (req, res) => {
  logs.buscar_logs()
  res.send('lido')
})

app.get('/', async (req, res) => {
  await logs.buscar_logs()
  res.json(db.data)
})

app.get('/apagar', (req, res) => {
  db.data.length = 0
  res.send('Apagado')
})

app.listen(port, () => {
  console.log(`Example app listening on port ${port}`)
})