# scrapad_backend
 
He modifcado la base de datos añadiendo el campo ad_id en la tabla offers para relacionar cada oferta con el anuncio desde el que se ha realizado. Esto es necesario ya que se requiere que al devolver las ofertas aparezca el anuncio desde el que se han realizado, y no hay forma de saber esto con los campos y relaciones actuales. Ademas, al crear la oferta tambien es necesario el id del anuncio desde el que se realiza, por lo que este nuevo campo es necesario

He modificado las fechas de creacion de las orgs, añadiéndoles un año a todas, dado que ninguna actualmente tiene un año de antigüedad a dia de hoy.