# frozen_string_literal: true

require 'nokogiri'
require 'time'
require 'uri'
require 'cgi'
require 'json'

name = ARGV[0]

emojis = {
  '<i class="_emReg" data-emoji=":)"></i>' => '😊',
  '<i class="_emReg" data-emoji=":P"></i>' => '😛',
  '<i class="_emReg" data-emoji="😊"></i>' => '😊',
  '<i class="_emReg" data-emoji="(y)"></i>' => '👍',
  '<i class="_emReg" data-emoji=":D"></i>' => '😃',
  '<i class="_emReg" data-emoji="😛"></i>' => '😛',
  '<i class="_emReg" data-emoji="😕"></i>' => '😕',
  '<i class="_emReg" data-emoji="😉"></i>' => '😉',
  '<i class="_emReg" data-emoji="😂"></i>' => '😂',
  '<i class="_emReg" data-emoji=":p"></i>' => '😛',
  '<i class="_emReg" data-emoji="👌"></i>' => '👌',
  '<i class="_emReg" data-emoji=";)"></i>' => '😉',
  '<i class="_emReg" data-emoji="😀"></i>' => '😀',
  '<i class="_emReg" data-emoji="🙂"></i>' => '🙂',
  '<i class="_emReg" data-emoji=":o"></i>' => '😮',
  '<i class="_emReg" data-emoji=":("></i>' => '😟',
  '<i class="_emReg" data-emoji="😴"></i>' => '😴',
  '<i class="_emReg" data-emoji="😲"></i>' => '😲',
  '<i class="_emReg" data-emoji="💕"></i>' => '💕',
  '<i class="_emReg" data-emoji="😭"></i>' => '😭',
  '<i class="_emReg" data-emoji="❤️"></i>' => '❤️',
  '<i class="_emReg" data-emoji="<3"></i>' => '❤️',
  '<i class="_emReg" data-emoji="❤"></i>' => '❤️',
  '<i class="_emReg" data-emoji="😁"></i>' => '😁',
  '<i class="_emReg" data-emoji="😗"></i>' => '😗',
  '<i class="_emReg" data-emoji="🍆"></i>' => '🍆',
  '<i class="_emReg" data-emoji="😃"></i>' => '😃',
  '<i class="_emReg" data-emoji="👍"></i>' => '👍',
  '<i class="_emReg" data-emoji="🇬🇧"></i>' => '🇬🇧',
  '<i class="_emReg" data-emoji="😔"></i>' => '😔',
  '<i class="_emReg" data-emoji=":/"></i>' => '😕',
  '<i class="_emReg" data-emoji="😍"></i>' => '😍',
  '<i class="_emReg" data-emoji="🎶"></i>' => '🎶',
  '<i class="_emReg" data-emoji="✌️"></i>' => '✌️',
  '<i class="_emReg" data-emoji="😄"></i>' => '😄',
  '<i class="_emReg" data-emoji="😅"></i>' => '😅',
  '<i class="_emReg" data-emoji="😎"></i>' => '😎',
  '<i class="_emReg" data-emoji="🔫"></i>' => '🔫',
  '<i class="_emReg" data-emoji="🗿"></i>' => '🗿',
  '<i class="_emReg" data-emoji="😇"></i>' => '😇',
  '<i class="_emReg" data-emoji="🅱"></i>' => '🅱',
  '<i class="_emReg" data-emoji="🅰"></i>' => '🅰',
  '<i class="_emReg" data-emoji="🅾"></i>' => '🅾',
  '<i class="_emReg" data-emoji="🆎"></i>' => '🆎',
  '<i class="_emReg" data-emoji="🆖"></i>' => '🆖',
  '<i class="_emReg" data-emoji="🔞"></i>' => '🔞',
  '<i class="_emReg" data-emoji="🏵"></i>' => '🏵',
  '<i class="_emReg" data-emoji="🍑"></i>' => '🍑',
  '<i class="_emReg" data-emoji="👈"></i>' => '👈',
  '<i class="_emReg" data-emoji="😱"></i>' => '😱',
  '<i class="_emReg" data-emoji="😘"></i>' => '😘',
  '<i class="_emReg" data-emoji="😵"></i>' => '😵',
  '<i class="_emReg" data-emoji="😮"></i>' => '😮',
  '<i class="_emReg" data-emoji="🎉"></i>' => '🎉',
  '<i class="_emReg" data-emoji="😬"></i>' => '😬',
  '<i class="_emReg" data-emoji="🤑"></i>' => '🤑',
  '<i class="_emReg" data-emoji="=D"></i>' => '=D',
  '<i class="_emReg" data-emoji="💣"></i>' => '💣',
  '<i class="_emReg" data-emoji="😤"></i>' => '😤',
  '<i class="_emReg" data-emoji="😞"></i>' => '😞',
  '<i class="_emReg" data-emoji="☺"></i>' => '☺',
  '<i class="_emReg" data-emoji=":-)"></i>' => ':-)',
  '<i class="_emReg" data-emoji=";-)"></i>' => ';-)',
  '<i class="_emReg" data-emoji=":3"></i>' => ':3',
  '<i class="_emReg" data-emoji=":\'("></i>' => ':\'(',
  '<i class="_emReg" data-emoji="🎯"></i>' => '🎯',
  '<i class="_emReg" data-emoji="🖕🏻"></i>' => '🖕🏻',
  '<i class="_emReg" data-emoji="😻"></i>' => '😻',
  '<i class="_emReg" data-emoji="😸"></i>' => '😸',
  '<i class="_emReg" data-emoji="😳"></i>' => '😳',
  '<i class="_emReg" data-emoji="☝🏽"></i>' => '☝🏽',
  '<i class="_emReg" data-emoji="😽"></i>' => '😽',
  '<i class="_emReg" data-emoji="😷"></i>' => '😷',
  '<i class="_emReg" data-emoji="😏"></i>' => '😏',
  '<i class="_emReg" data-emoji="😩"></i>' => '😩',
  '<i class="_emReg" data-emoji="😝"></i>' => '😝',
  '<i class="_emReg" data-emoji="🚀"></i>' => '🚀',
  '<i class="_emReg" data-emoji="(Y)"></i>' => '(Y)',
  '<i class="_emReg" data-emoji="🤔"></i>' => '🤔',
  '<i class="_emReg" data-emoji="🍻"></i>' => '🍻',
  '<i class="_emReg" data-emoji="🍺"></i>' => '🍺',
  '<i class="_emReg" data-emoji="👍🏿"></i>' => '👍🏿',
  '<i class="_emReg" data-emoji="🙊"></i>' => '🙊',
  '<i class="_emReg" data-emoji="😚"></i>' => '😚',
  '<i class="_emReg" data-emoji="💅🏿"></i>' => '💅🏿',
  '<i class="_emReg" data-emoji="🍎"></i>' => '🍎',
  '<i class="_emReg" data-emoji="👃"></i>' => '👃',
  '<i class="_emReg" data-emoji="👊"></i>' => '👊',
  '<i class="_emReg" data-emoji="&lt;3"></i>' => '&lt;3',
  '<i class="_emReg" data-emoji="🐬"></i>' => '🐬',
  '<i class="_emReg" data-emoji="🌈"></i>' => '🌈',
  '<i class="_emReg" data-emoji="😶"></i>' => '😶',
  '<i class="_emReg" data-emoji="🙏"></i>' => '🙏',
  '<i class="_emReg" data-emoji="👆"></i>' => '👆',
  '<i class="_emReg" data-emoji="🙋"></i>' => '🙋',
  '<i class="_emReg" data-emoji="☝"></i>' => '☝',
  '<i class="_emReg" data-emoji="✊"></i>' => '✊',
  '<i class="_emReg" data-emoji="🙃"></i>' => '🙃',
  '<i class="_emReg" data-emoji="👊🏾"></i>' => '👊🏾',
  '<i class="_emReg" data-emoji="💸"></i>' => '💸',
  '<i class="_emReg" data-emoji="😆"></i>' => '😆',
  '<i class="_emReg" data-emoji="👊🏼"></i>' => '👊🏼',
  '<i class="_emReg" data-emoji="🎻"></i>' => '🎻',
  '<i class="_emReg" data-emoji="😡"></i>' => '😡',
  '<i class="_emReg" data-emoji="😠"></i>' => '😠',
  '<i class="_emReg" data-emoji="-_-"></i>' => '-_-',
  '<i class="_emReg" data-emoji="☹️"></i>' => '☹️',
  '<i class="_emReg" data-emoji="🔥"></i>' => '🔥',
  '<i class="_emReg" data-emoji=":*"></i>' => ':*',
  '<i class="_emReg" data-emoji="👆🏼"></i>' => '👆🏼',
  '<i class="_emReg" data-emoji="😪"></i>' => '😪',
  '<i class="_emReg" data-emoji="😐"></i>' => '😐',
  '<i class="_emReg" data-emoji="🌸"></i>' => '🌸',
  '<i class="_emReg" data-emoji="👻"></i>' => '👻',
  '<i class="_emReg" data-emoji="👌🏿"></i>' => '👌🏿',
  '<i class="_emReg" data-emoji="😒"></i>' => '😒',
  '<i class="_emReg" data-emoji="🐧"></i>' => '🐧',
  '<i class="_emReg" data-emoji="🖖🏼"></i>' => '🖖🏼',
  '<i class="_emReg" data-emoji="👉"></i>' => '👉',
  '<i class="_emReg" data-emoji="👅"></i>' => '👅',
  '<i class="_emReg" data-emoji="🖕"></i>' => '🖕',
  '<i class="_emReg" data-emoji="💋"></i>' => '💋',
  '<i class="_emReg" data-emoji=":O"></i>' => ':O',
  '<i class="_emReg" data-emoji="🎩"></i>' => '🎩',
  '<i class="_emReg" data-emoji="💦"></i>' => '💦',
  '<i class="_emReg" data-emoji="😌"></i>' => '😌',
  '<i class="_emReg" data-emoji="🖕🏽"></i>' => '🖕🏽',
  '<i class="_emReg" data-emoji="⚽️"></i>' => '⚽️',
  '<i class="_emReg" data-emoji="👌🏻"></i>' => '👌🏻',
  '<i class="_emReg" data-emoji="🌩"></i>' => '🌩',
  '<i class="_emReg" data-emoji="😥"></i>' => '😥',
  '<i class="_emReg" data-emoji="💪"></i>' => '💪',
  '<i class="_emReg" data-emoji="✋🏻"></i>' => '✋🏻',
  '<i class="_emReg" data-emoji="💃🏿"></i>' => '💃🏿',
  '<i class="_emReg" data-emoji="🎊"></i>' => '🎊',
  '<i class="_emReg" data-emoji="🎆"></i>' => '🎆',
  '<i class="_emReg" data-emoji="🎇"></i>' => '🎇',
  '<i class="_emReg" data-emoji="💩"></i>' => '💩',
  '<i class="_emReg" data-emoji="☀"></i>' => '☀',
  '<i class="_emReg" data-emoji="😋"></i>' => '😋',
  '<i class="_emReg" data-emoji="🇩🇪"></i>' => '🇩🇪',
  '<i class="_emReg" data-emoji="🇫🇷"></i>' => '🇫🇷',
  '<i class="_emReg" data-emoji="🐻"></i>' => '🐻',
  '<i class="_emReg" data-emoji="😖"></i>' => '😖',
  '<i class="_emReg" data-emoji="😙"></i>' => '😙',
  '<i class="_emReg" data-emoji="✌"></i>' => '✌',
  '<i class="_emReg" data-emoji="😰"></i>' => '😰',
  '<i class="_emReg" data-emoji="💞"></i>' => '💞',
  '<i class="_emReg" data-emoji="🤘🏻"></i>' => '🤘🏻',
  '<i class="_emReg" data-emoji="🖕🏿"></i>' => '🖕🏿',
  '<i class="_emReg" data-emoji="✈"></i>' => '✈',
  '<i class="_emReg" data-emoji="😑"></i>' => '😑',
  '<i class="_emReg" data-emoji="😫"></i>' => '😫',
  '<i class="_emReg" data-emoji="🍟"></i>' => '🍟',
  '<i class="_emReg" data-emoji="🏀"></i>' => '🏀',
  '<i class="_emReg" data-emoji="🍅"></i>' => '🍅',
  '<i class="_emReg" data-emoji="🇮🇪"></i>' => '🇮🇪',
  '<i class="_emReg" data-emoji="✋"></i>' => '✋',
  '<i class="_emReg" data-emoji="☀️"></i>' => '☀️',
  '<i class="_emReg" data-emoji="👙"></i>' => '👙',
  '<i class="_emReg" data-emoji="✋🏾"></i>' => '✋🏾"',
  '<i class="_emReg" data-emoji="🖐🏻"></i>' => '🖐🏻',
  '<i class="_emReg" data-emoji="🖐🏼"></i>' => '🖐🏼',
  '<i class="_emReg" data-emoji="🖐🏽"></i>' => '🖐🏽',
  '<i class="_emReg" data-emoji="🖐🏾"></i>' => '🖐🏾',
  '<i class="_emReg" data-emoji="🖐🏿"></i>' => '🖐🏿',
  '<i class="_emReg" data-emoji="🖐"></i>' => '🖐',
  '<i class="_emReg" data-emoji="👫"></i>' => '👫',
  '<i class="_emReg" data-emoji="🐤"></i>' => '🐤',
  '<i class="_emReg" data-emoji="✋🏽"></i>' => '✋🏽',
  '<i class="_emReg" data-emoji="🎁"></i>' => '🎁',
  '<i class="_emReg" data-emoji="🐶"></i>' => '🐶',
  '<i class="_emReg" data-emoji="✋🏼"></i>' => '✋🏼',
  '<i class="_emReg" data-emoji="💍"></i>' => '💍',
  '<i class="_emReg" data-emoji="😓"></i>' => '😓',
  '<i class="_emReg" data-emoji="🙏🏻"></i>' => '🙏🏻',
  '<i class="_emReg" data-emoji="🍣"></i>' => '🍣',
  '<i class="_emReg" data-emoji="✋🏿"></i>' => '✋🏿',
  '<i class="_emReg" data-emoji="🐼"></i>' => '🐼',
  '<i class="_emReg" data-emoji="👀"></i>' => '👀',
  '<i class="_emReg" data-emoji="🎱"></i>' => '🎱',
  '<i class="_emReg" data-emoji="🙈"></i>' => '🙈',
  '<i class="_emReg" data-emoji="👊🏿"></i>' => '👊🏿',
  '<i class="_emReg" data-emoji="✅"></i>' => '✅',
  '<i class="_emReg" data-emoji="🎄"></i>' => '🎄',
  '<i class="_emReg" data-emoji="🍭"></i>' => '🍭',
}

html = File.read("/Users/ak/Downloads/#{name}/#{name}.html")
html = Nokogiri::HTML(html).to_html

emojis.each do |code, replacement|
  html.gsub!(code, replacement)
end

doc = Nokogiri::HTML(html, nil, 'UTF-8')

found_emojis = doc.css('._emReg')
unless found_emojis.empty?
  puts found_emojis.map(&:to_html)
  raise 'unreplaced emojis'
end

results = []

doc.css('._42ef').each do |msg_div|
  unless msg_div.css('._1xpw').empty?
    # notice message
    next
  end

  timestamp = msg_div.css('.timestamp').first['data-utime'].to_i / 1000
  author = msg_div.css('.aClass').first.content

  if author == 'Adrien Kohlbecker'
    sent = true
    author = ''
  else
    sent = false
  end

  result_tmpl = {
    'id' => '',
    'sender' => author,
    'content' => '',
    'timestamp' => Time.at(timestamp).iso8601,
    'sent' => sent,
    'attachments' => [],
    'group' => name,
    'kind' => 'facebook',
    'replies_to' => nil
  }

  msg_div.children.each do |msg_div_child|

    if msg_div_child['class'].nil?

        msg_div_child.children.each do |div|

            if div['class'] == '_36' # author
                next
            elsif div['class'] == '_37' # messages
                div.css('.pClass').each do |msg|

                    content = msg.inner_html.strip
                    content.gsub!('<br>', "\n")
                    content.gsub!(/<a[^>]*>/, '')
                    content.gsub!('</a>', '')

                    result = result_tmpl.dup
                    result['timestamp'] = Time.at(timestamp).iso8601
                    result['content'] = Nokogiri::HTML(content).content
                    results << result
                    timestamp = timestamp + 1

                end

            elsif div['class'] == '_sq' # link

                if div.css('._3bwx').length > 0 # video attachement

                    result = result_tmpl.dup
                    result['timestamp'] = Time.at(timestamp).iso8601
                    result['attachments'] = [{
                      'kind' => 'video',
                      'url' => "/Users/ak/Google Drive/Applications/Messages/media/" + div.css('a').first['href']
                    }]
                    results << result
                    timestamp = timestamp + 1

                elsif div.css('._59gp').length > 0 # document attachement

                    result = result_tmpl.dup
                    result['timestamp'] = Time.at(timestamp).iso8601
                    result['attachments'] = [{
                      'kind' => 'document',
                      'url' => "/Users/ak/Google Drive/Applications/Messages/media/" + div.css('._59gp').first.content
                    }]
                    results << result
                    timestamp = timestamp + 1

                else # plain link

                    link = div.css('a._5rw4').first['href']
                    if link.include?('https://l.facebook.com/l.php')
                        # https://l.facebook.com/l.php?u=https%3A%2F%2Fwww.residentadvisor.net%2Fevent.aspx%3F808746&h=VAQExXFWz&s=1
                        link = CGI::parse(URI.parse(link).query)['u'][0]
                    end

                    result = result_tmpl.dup
                    result['timestamp'] = Time.at(timestamp).iso8601
                    result['content'] = link
                    results << result
                    timestamp = timestamp + 1

                end

            elsif div.name == 'img' # image attachment

                result = result_tmpl.dup
                result['timestamp'] = Time.at(timestamp).iso8601
                result['attachments'] = [{
                  'kind' => 'img',
                  'url' => "/Users/ak/Google Drive/Applications/Messages/media/" + div['src']
                }]
                results << result
                timestamp = timestamp + 1

            elsif div.text?
                next
            else
                puts msg_div.to_html
                raise div.inspect
            end

        end

    end
  end

end

File.write("#{name}.json", JSON.pretty_generate(results))
