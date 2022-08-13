const fs = require('fs')
const readline = require('readline')
const base = require('./logs')

const dir = './logs/'

const logs = {
  id: 0,
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
    logs.id = base.id
    base.id++
    base.data.push({ ...logs })
    remover_arquivo(arq)
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
  logs.id     = 0
  logs.metodo = ''
  logs.horas  = ''
  logs.sql    = []
  logs.outros = []
  logs.requisicao = ''
  logs.resposta   = ''
}

const retornaLogs = id => {
  return base.data.filter(item => item.id > id)
}

const ligaDesliga = async value => {
  const arq = await fs.readFileSync(dir + 'config.conf').toString('UTF8').split('\n')
  
  if (arq[1] !== value && value !== -1) {
    arq[1] = value
    const linhas = arq[0] + '\n' + arq[1]
    fs.writeFileSync(dir + 'config.conf', linhas)
  }

  return (arq[1] === '0') ? 'Off' : 'On'
}

module.exports = {
  buscar_logs,
  ler_log,
  retornaLogs,
  ligaDesliga
}