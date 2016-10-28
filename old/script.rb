require 'base64'

file = File.read('heliot-preprocessing.xml')
out = File.open('heliot-post.xml', 'w')
out.write(
  <<-MMS
  <?xml version='1.0' encoding='UTF-8' standalone='yes' ?>
  <!--File Created By SMS Backup & Restore v8.20.27 on 21/07/2016 13:15:13-->
  <?xml-stylesheet type="text/xsl" href="sms.xsl"?>
  <smses count="15" backup_set="bdd104eb-f636-4fae-aad9-e3ce37a649ac" backup_date="1469099713413">
  MMS
)

file.scan(%r{<sms.*address="(?<address>[^"]*)".*date="(?<date>[^"]*)".*body="(?<body>[^"]*)".*/>}).each do |match|
  address = match[0]
  date = match[1]
  body = match[2]

  address_tilde = '+33629448767~+33695218383'

  if address == '+33629448767'
    addrs = <<-ADDR
      <addr address="+33695218383" type="151" charset="106" />
      <addr address="+33629448767" type="137" charset="106" />
    ADDR
    msg_box = '1'
  elsif address == '+33 6 95 21 83 83'
    addrs = <<-ADDR
      <addr address="+33695218383" type="137" charset="106" />
      <addr address="+33629448767" type="151" charset="106" />
    ADDR
    msg_box = '1'
  else
    addrs = <<-ADDR
      <addr address="+33695218383" type="151" charset="106" />
      <addr address="+33629448767" type="151" charset="106" />
    ADDR
    msg_box = '2'
  end

  if body.start_with?('/Users')

    if body.end_with?('.jpeg', '.jpg')
      type = 'image/jpeg'
    elsif body.end_with?('.3gp')
      type = 'video/3gpp'
    end

    encoded = Base64.encode64(File.read(body)).strip.delete("\n")
    filename = body.split('/').last
    name = filename.split('.').first
    part = "<part seq=\"0\" ct=\"#{type}\" name=\"null\" chset=\"null\" cd=\"null\" fn=\"null\" cid=\"&lt;#{name}&gt;\" cl=\"#{filename}\" ctt_s=\"null\" ctt_t=\"null\" text=\"null\" data=\"#{encoded}\" />"

  else
    part = "<part seq=\"0\" ct=\"text/plain\" name=\"part-0\" chset=\"106\" cd=\"null\" fn=\"part-0\" cid=\"null\" cl=\"null\" ctt_s=\"null\" ctt_t=\"null\" text=\"#{body}\" />"
  end

  out.write(
    <<-MMS
    <mms text_only="0" group_id="null" ct_t="application/vnd.wap.multipart.related" reserve_time="0" msg_box="#{msg_box}" v="18" ct_cls="null" retr_txt_cs="null" phone_id="-1" type="-1" st="null" save_call_type="null" creator="com.textra" tag="null" tr_id="null" read="1" m_id="null" m_type="128" name="null" locked="0" retr_txt="null" resp_txt="null" retr_st="null" sub="null" seen="0" rr="129" ct_l="null" m_size="null" exp="null" c0_iei="0" sub_cs="null" sub_id="-1" imsi_data="null" resp_st="null" date="#{date}" date_sent="0" pri="129" msg_boxtype="null" textlink="-1" insert_time="#{date}" address="#{address_tilde}" d_rpt="129" d_tm="null" read_status="null" spam_report="0" sender_num="null" rpt_a="null" kt_tm_send_type="null" m_cls="personal" tag_eng="null">
      <parts>
        #{part}
      </parts>
      <addrs>#{addrs}</addrs>
    </mms>
    MMS
  )
end

out.write('</smses>')
