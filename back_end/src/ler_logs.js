const fs = require('fs')
const readline = require('readline')
const { data } = require('./logs')

const dir = './logs/'

const logs = {
  metodo: '',
  horas: '',
  sql: [],
  outros: [],
  requisicao: '',
  resposta: ''
}

const buscar_logs = async () => {
  const lista = await fs.readdirSync(dir)
  for (const arq of lista) {
    if (arq === 'config.conf')
     continue

    await ler_arquivo(arq)
    data.push({ ...logs })
    //remover_arquivo(arq)
    limpar()
  }
}

const ler_log = line => {
  switch (line.substr(0, 3)) {
    case 'met':
      logs.metodo = line.substr(5)
      break
    case 'sql':
      logs.sql.push(line.substr(5))
      break
    case 'req':
      logs.requisicao = line.substr(5)
      break
    case 'res':
      logs.resposta = line.substr(5)
      break
    case 'hor':
      logs.horas = line.substr(5)
      break
    default:
      logs.outros.push(line)
      break
  }

}

const ler_arquivo = async path => {
  const arq = fs.createReadStream(dir + path)

  const rl = readline.createInterface({
    input: arq,
    crlfDelay: Infinity
  })
  
  for await (const line of rl) {
    ler_log(line)
  }
} 

const remover_arquivo = path => {
  fs.unlink(dir + path, (err => {
    if (err) console.log(err);
    else {
      console.log("\nDeleted file: " + path);
    }
  }));
}

const limpar = () => {
  logs.metodo = ''
  logs.horas  = ''
  logs.sql    = []
  logs.outros = []
  logs.requisicao = ''
  logs.resposta   = ''
}

module.exports = {
  buscar_logs,
  ler_log
}