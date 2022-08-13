const copiarTexto = () => {
  let textoCopiado = 'opas'
  navigator.clipboard.writeText(textoCopiado)
}

const url = 'http://localhost:5000/'
const baseDados = [];
const config = { filter: false }

const loadLogsAPI = () => {
  const point = ((baseDados.length > 0) ? 'ler/' + baseDados[0].id + '/' + baseDados[0].horas : 'ler')

  $.ajax({
    url: url + point, 
    success: data => {
      loadLogs(data)
    },
    error: data => {
      console.log(data)
    }
  })
}

const limparLogsAPI = () => {
  $.ajax({
    url: url + "apagar", 
    success: data => {
      baseDados.length = 0
      drawLogs()
      $('.info').css("display", "none")
      $('.vazio').css("display", "flex")      
      console.log('Logs Limpados')
    },
    error: data => {
      console.log(data)
    }
  })
}

const pesquisarLogs = () => {
  const text = $('.pesquisa').val()

  if (text.length > 0) {
    const logsFilter = baseDados.filter(log => {
      return log.metodo.includes(text)
    })

    config.filter = true
    drawLogs((logsFilter.length > 0 ? logsFilter : []))
  } else if ( config.filter ) {
    drawLogs()
  }
}

const loadLogs = (api) => {
  api.map(log => {
    log.sql.map((sql, i) => {
      log.sql[i] = sqlFormat(sql)
    })
    log.requisicao = xmlFormat(log.requisicao)
    log.resposta = xmlFormat(log.resposta)

    log.status = 0

    baseDados.unshift(log)
  })
  drawLogs()
}

const drawLogs = (logs = null) => {
  let htmlText = ''
  if (logs === null)
    logs = baseDados

  logs.map((log, i) => {
    htmlText += `<div class="row log ${(log.status === 0) ? "border-info" : ""} align-items-center justify-content-center"` +
              `  onclick="loadInfo(${i})" > ` +
              '   <div class="row">' +
              `     <div class="col-6">${log.metodo}</div> ` +
              `     <div class="col-6 text-center">${log.horas}</div>` +
              '   </div> ' +
              ' </div>'

    log.status = 1              

  })  

  if (logs.length === 0) {
    htmlText = '<div class="d-flex h-100 justify-content-center align-items-center"><span>Nenhum log encontrado!</span></div>'
  }
      
  $('.list-logs').html(htmlText)
}

const loadInfo = index => {
  const log = baseDados[index]

  drawSQL(log.sql)
  drawOut(log.outros)
  drawReq(log.requisicao)
  drawRes(log.resposta)

  $('.info').css("display", "block")
  $('.vazio').css("display", "none")
}

const drawSQL = sqls => {
  let htmlText = ''

  sqls.map(sql => {
    htmlText += `<textarea class="log_format" rows="${sql.line}" disabled>${sql.sql}</textarea>`
  })

  $('.log-sql').html(htmlText)
}

const drawOut = out => {
  let htmlText = ''

  out.map(iten => {
    htmlText += `<textarea class="log_format" rows="1" disabled>${iten}</textarea>`
  })

  $('.log-out').html(htmlText)
}

const drawReq = req => {
  const htmlText = `<textarea class="log_format" rows="${req.line}" disabled>${req.xml}</textarea>`
  
  $('.log-req').html(htmlText)
}

const drawRes = res => {
  const htmlText = `<textarea class="log_format" rows="${res.line}" disabled>${res.xml}</textarea>`
  
  $('.log-res').html(htmlText)
}

$(window).on("load", function(){
  loadLogsAPI()
});