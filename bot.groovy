/*!
 *  Crawler
 */
import java.net.URLEncoder
import groovy.time.TimeCategory
import groovy.json.JsonSlurper
import groovy.util.XmlSlurper
import groovyx.gpars.GParsPool
import sun.misc.BASE64Decoder
import org.jsoup.Jsoup
import org.jsoup.Connection.Method


ebook = [
    '1': ['莽荒记', '3961103225'],
    '2': ['网游之狂兽逆天', '1954128079'],
    '3': ['蜗居', '1822469960'],
    '4': ['异世灵武天下', '1705173889'],
    '5': ['全职高手', '4233102889']
]

// default action
gid = ebook['1'][1]
if (this.args.length) {
    gid = ebook[this.args[0]][1]
}




debug = 1
if (debug == 1) {
    println "setting proxy -> proxy.pvgl.sap.corp:8080"
    System.setProperty("http.proxyHost", "proxy.pvgl.sap.corp")
    System.setProperty("http.proxyPort", "8080")
    System.setProperty("https.proxyHost", "proxy.pvgl.sap.corp")
    System.setProperty("https.proxyPort", "8080")
} else if (debug == 2) {
    println "setting proxy -> localhost:8888"
    System.setProperty("http.proxyHost", "localhost")
    System.setProperty("http.proxyPort", "8888")
    System.setProperty("https.proxyHost", "localhost")
    System.setProperty("https.proxyPort", "8888")
}


DEFAULT_TIMEOUT = 30000

def fetchContents() {
    def action = "http://m.baidu.com/tc?srd=1&appui=alaxs&ajax=1&gid=${gid}"
    def content = new JsonSlurper().parse(new URL(action), "utf-8")
    return content;
}

def fetchArticles(content) {
    def slurper = new JsonSlurper()
    return GParsPool.withPool(32) {
        content.data.group.collectParallel {
            println " retreiving ${it.href}"
            def cid = URLEncoder.encode(it.cid)
            def href = URLEncoder.encode(it.href)
            // def action = "http://m.baidu.com/tc?srd=1&appui=alaxs&ajax=1&gid=1954128079&alals=1&preNum=1&preLoad=true&src=${href}&cid=${cid}&time=&skey=&id=wisenovel"
            def action = "http://m.baidu.com/tc?srd=1&appui=alaxs&ajax=1&gid=${gid}&alals=1&preNum=1&preLoad=true&src=${href}&cid=${cid}&time=&skey=&id=wisenovel"
            //println " cc..."
            def article = slurper.parse(new URL(action), "utf-8")
            //println " cc: ${article}"
            return [href: it.href, title: it.text, article: article]
        }
    }.each {
        println " processing ${it.href}"
        //it.content = slurper.parse(it.article, "utf-8").content
        try {
            def cc = it.article.data.first().content.replace('<br/>','__BR__')
            cc = Jsoup.parse(cc).text().replace('__BR__', '\n')
            it.content = cc
        }
        catch(Exception e) {
            println " failed to process ${it.article}"
            it.content = "BAD_CONTENT"
        }
    }
}

def save(content, articles) {
    new File(content.data.title + '.txt').withWriter("unicode") { fs ->
        fs.println "[${content.data.title}]"
        fs.println "作者: ${content.data.author}"
        fs.println ""

        articles.each {
            fs.println "[${it.title}]"
            fs.println it.content
            fs.println ""
        }
    }
}

def main() {
    def start = new Date()

    println "crawler is working..."

    def content = fetchContents()
    println "fetch content, src=${content.data.src}"

    def articles = fetchArticles(content)
    println "fetch articles successfully"

    save(content, articles)
    println "save ebook successfully"

    def end = new Date()
    def elapsed = TimeCategory.minus(end, start)
    println "all done, time elapsed: ${elapsed}"
}

main()
