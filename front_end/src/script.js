

const url = 'http://localhost:5000/'
const baseDados = [];
const config = { filter: false }

function ativarDesativar(){
  const check = $('#ckb-on-off')[0].checked

  $.ajax({
    url: url + 'status/' + ((check) ? 0 : 1),
    success: data => {
      console.log(data)
    },
    error: data => {
      console.log(data)
    }
  })  
}

function statusLogs(){
  $.ajax({
    url: url + 'status/',
    success: data => {
      $('#ckb-on-off')[0].checked = (data == '0')
    },
    error: data => {
      console.log(data)
    }
  })  
}
setInterval(statusLogs, 60000);

const loadLogsAPI = () => {
  let higherLog
  if (baseDados.length > 0) {
    higherLog = baseDados.reduce(function(prev, current) {
      return (prev.Id > current.Id) ? prev : current
    })
  }
  const point = ((higherLog) ? 'log/' + higherLog.Id + '/' + higherLog.Horas : 'log/')

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
    url: url + "limpar", 
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
      return log.Metodo.includes(text)
    })

    config.filter = true
    drawLogs((logsFilter.length > 0 ? logsFilter : []))
  } else if ( config.filter ) {
    drawLogs()
  }
}

const loadLogs = (api) => {
  api.map(log => {
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
              `  onclick="loadInfo(${log.Id})" > ` +
              '   <div class="row">' +
              `     <div class="col-6">${log.Metodo}</div> ` +
              `     <div class="col-6 text-center">${log.Horas}</div>` +
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
  const log = baseDados.filter(log => log.Id == index)[0]
  $('#titleMetodo').text(log.Metodo)

  drawSQL(log.Sql)
  //drawOut(log.Outros)
  drawReq(log.Requisicao)
  drawRes(log.Resposta)

  $('.info').css("display", "block")
  $('.vazio').css("display", "none")
}

const drawSQL = sqls => {
  let htmlText = ''

  sqls.map(sql => {
    htmlText += `<textarea class="log_format" rows="${sql.Linhas}" wrap="off" disabled>${sql.Script}</textarea>`
  })

  $('.log-sql').html(htmlText)
}

const drawOut = out => {
  let htmlText = ''

  out.map(iten => {
    htmlText += `<textarea class="log_format" rows="1" wrap="off" disabled>${iten}</textarea>`
  })

  $('.log-out').html(htmlText)
}

const drawReq = req => {
  const htmlText = `<textarea class="log_format" rows="${req.Linhas}" wrap="off" disabled>${req.Script}</textarea>`
  
  $('.log-req').html(htmlText)
}

const drawRes = res => {
  const htmlText = `<textarea class="log_format" rows="${res.Linhas}" wrap="off" disabled>${res.Script}</textarea>`
  
  $('.log-res').html(htmlText)
}

$(window).on("load", function(){
  statusLogs()
  loadLogsAPI()
});