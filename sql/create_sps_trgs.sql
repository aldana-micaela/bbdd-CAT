--SPs

create or replace function ingresar_llamado(dato_id_cliente int) returns int as $$
declare
    v_id_cola_atencion int;
    v_cliente_existe boolean;
begin
    select exists (select 1 from cliente where id_cliente = dato_id_cliente) into v_cliente_existe;
    
	if not v_cliente_existe then
        insert into error values (default, 'nuevo llamado', null, null, null, null, null, CURRENT_TIMESTAMP(0)::TIMESTAMP, 'id de cliente no válido');
        return -1; 
    end if;
   
    insert into cola_atencion values (default, dato_id_cliente, CURRENT_TIMESTAMP(0)::TIMESTAMP, null, null, null, 'en espera')
    returning id_cola_atencion into v_id_cola_atencion;
    
    return v_id_cola_atencion;
end;
$$ language plpgsql;



create or replace function desistir_llamado (id_cola_at int) returns boolean as $$
declare
	llamado_desistido boolean;
	id_existe boolean;
	estado_llamado char (15);
    v_id_operadore int;
    v_id_cliente int;
begin
	
	select exists (select 1 from cola_atencion where id_cola_atencion = id_cola_at) into id_existe;
	if not id_existe then
		insert into error values (default, 'baja llamado', null, null, null, null, null, CURRENT_TIMESTAMP(0)::TIMESTAMP, '?id de cola de atención no válido');
		return false;
	end if;

	select estado, id_cliente into estado_llamado, v_id_cliente from cola_atencion where id_cola_atencion = id_cola_at;
	if estado_llamado !='en espera' and estado_llamado != 'en linea' then 
		insert into error values (default, 'baja llamado', v_id_cliente, id_cola_at, null, null, null, CURRENT_TIMESTAMP(0)::TIMESTAMP, '?el llamado no está en espera ni en línea');
		return false;
	end if;

	 if estado_llamado = 'en linea' then
        update cola_atencion set estado = 'desistido', f_fin_atencion = CURRENT_TIMESTAMP(0)::TIMESTAMP where id_cola_atencion = id_cola_at;
        
        select id_operadore into v_id_operadore from cola_atencion where id_cola_atencion = id_cola_at;
	    update operadore set disponible = true where id_operadore= v_id_operadore;


    else
        update cola_atencion set estado = 'desistido' where id_cola_atencion = id_cola_at;
    end if;

	return true;
end;
$$ language plpgsql;


create or replace function atender_llamado() returns boolean as $$
declare
    v_id_cola_atencion int;
    v_id_operadore int;
    v_id_cliente int;
    v_llamado_existe boolean;
    operadore_disponible boolean;
begin
    select exists (select 1 from cola_atencion where estado = 'en espera')
    into v_llamado_existe;

    if not v_llamado_existe then
        insert into error values (default, 'atencion llamado', null, null, null, null, null, CURRENT_TIMESTAMP(0)::TIMESTAMP, '?no existe ningún llamado en espera');
        return false;
    end if;

    select id_cliente into v_id_cliente from cola_atencion where estado = 'en espera'
    order by f_inicio_llamado
    limit 1;

    select exists (select disponible from operadore where disponible = true) into operadore_disponible;

    if not operadore_disponible then
        insert into error values (default, 'atencion llamado', v_id_cliente, null, null, null, null, CURRENT_TIMESTAMP(0)::TIMESTAMP, '?no existe ningun operadore disponible');
        return false;
    end if;

    select id_operadore into v_id_operadore from operadore where disponible = true
    order by fecha_ingreso
    limit 1;

    update operadore set disponible = false where id_operadore = v_id_operadore;

    select id_cola_atencion into v_id_cola_atencion from cola_atencion where estado = 'en espera'
    order by f_inicio_llamado
    limit 1;

    update cola_atencion set id_operadore = v_id_operadore,
        f_inicio_atencion = CURRENT_TIMESTAMP(0)::TIMESTAMP,
        estado = 'en linea'
    where id_cola_atencion = v_id_cola_atencion;

    return true;
end;
$$ language plpgsql;



create or replace function alta_de_tramite(id_cola_at int, t_tramite char(10), descrip text ) returns int as $$
declare
    v_id_tramite int;
    v_id_cola_at_existe boolean;
    v_estado char (15);
    v_id_cliente int;
begin

    if t_tramite != 'consulta' and t_tramite != 'reclamo' then
        insert into error values (default, 'alta tramite', null, null, null, null, null, CURRENT_TIMESTAMP(0)::TIMESTAMP, '?tipo de trámite no válido');
        return -1;
    end if;

    select exists (select 1 from cola_atencion where id_cola_atencion = id_cola_at ) into v_id_cola_at_existe;
    select estado into v_estado from cola_atencion where id_cola_atencion = id_cola_at;

    if not v_id_cola_at_existe or v_estado = 'en espera' then
        insert into error values (default, 'alta tramite', null, null, null, null, null, CURRENT_TIMESTAMP(0)::TIMESTAMP,'?id de cola de atención no válido');
       return -1;
    end if;

    select id_cliente into v_id_cliente from cola_atencion where id_cola_atencion = id_cola_at;

    insert into tramite values (default, v_id_cliente, id_cola_at, t_tramite, CURRENT_TIMESTAMP(0)::TIMESTAMP, descrip, null, null, 'iniciado')
    returning id_tramite into v_id_tramite;
    return v_id_tramite;

end;
$$ language plpgsql;

create or replace function finalizar_llamado(p_id_cola_atencion int) returns boolean as $$
declare
    v_estado char(15);
    v_id_cola_at_existe boolean;
	v_id_cliente int;
    v_id_operadore int;
begin
    select exists (select 1 from cola_atencion
	where id_cola_atencion = p_id_cola_atencion)
    into v_id_cola_at_existe;

    if not v_id_cola_at_existe then
        insert into error values (default, 'fin llamado', null, null, null, null, null,
		CURRENT_TIMESTAMP(0)::TIMESTAMP,'?id de cola de atención no válido');
        return false;
    end if;

    select estado into v_estado from cola_atencion
    where id_cola_atencion = p_id_cola_atencion;

	select id_cliente into v_id_cliente from cola_atencion 
	where id_cola_atencion = p_id_cola_atencion;

    if v_estado != 'en linea' then
        insert into error values (default, 'fin llamado', v_id_cliente, p_id_cola_atencion, null, 
            null, null, CURRENT_TIMESTAMP(0)::TIMESTAMP, '?el llamado no está en línea');
        return false;
    end if;

    update cola_atencion
    set estado = 'finalizado', f_fin_atencion = CURRENT_TIMESTAMP(0)::TIMESTAMP
    where id_cola_atencion = p_id_cola_atencion;


	select id_operadore into v_id_operadore from cola_atencion where id_cola_atencion = p_id_cola_atencion;
	update operadore set disponible = true where id_operadore= v_id_operadore;

    return true;
end;
$$ language plpgsql;


create or replace function cierre_tramite(p_id_tramite int, p_estado_cierre char(15), p_respuesta_tramite text) 
returns boolean AS $$
declare
    v_estado_tramite char(15);
begin
    if p_estado_cierre not in ('solucionado', 'rechazado') then
        insert into error (operacion, id_tramite, f_error, motivo)
        values ('cierre_tramite', p_id_tramite, current_timestamp, 'estado de cierre no válido');
        return false;
    end if;

    select estado into v_estado_tramite
    from tramite
    where id_tramite = p_id_tramite;

    if v_estado_tramite is null then
        insert into error (operacion, id_tramite, f_error, motivo)
        values ('cierre_tramite', p_id_tramite, current_timestamp, 'id de trámite no válido');
        return false;
    end if;


    if v_estado_tramite != 'iniciado' then
        insert into error (operacion, id_tramite, f_error, motivo)
        values ('cierre_tramite', p_id_tramite, current_timestamp, 'el trámite se encuentra cerrado');
        return false;
    end if;

    update tramite
    set estado = p_estado_cierre,
        f_fin_gestion = current_timestamp,
        respuesta = p_respuesta_tramite
    where id_tramite = p_id_tramite;

    return true;
end;
$$ language plpgsql;



create or replace function reporte_rendimiento() returns trigger as $$
declare
	duracion_llamada interval;
	v_fecha_atencion date;
begin
		if new.f_fin_atencion is not null and old.f_inicio_atencion is not null and old.id_operadore is not null then
        	duracion_llamada := new.f_fin_atencion - old.f_inicio_atencion;
        	v_fecha_atencion := new.f_fin_atencion;
			
			update rendimiento_operadore set
	        duracion_total_atenciones = duracion_total_atenciones + duracion_llamada,
	        cantidad_total_atenciones = cantidad_total_atenciones + 1,
	        duracion_promedio_total_atenciones = (duracion_total_atenciones + duracion_llamada) / (cantidad_total_atenciones + 1)
			where id_operadore = old.id_operadore and fecha_atencion = v_fecha_atencion;

			if not found then
				insert into rendimiento_operadore values (old.id_operadore, v_fecha_atencion, 
				duracion_llamada, 1, duracion_llamada, 
                '00:00:00', 0, '00:00:00',
				'00:00:00', 0, '00:00:00' );
			end if;

			if new.estado = 'desistido' then
				update rendimiento_operadore set
	            duracion_atenciones_desistidas = duracion_atenciones_desistidas + duracion_llamada,
	            cantidad_atenciones_desistidas = cantidad_atenciones_desistidas + 1,
	            duracion_promedio_atenciones_desistidas = (duracion_atenciones_desistidas + duracion_llamada) /(cantidad_atenciones_desistidas +1)
	       		where id_operadore = old.id_operadore and fecha_atencion = v_fecha_atencion;
			end if;
			
			if new.estado = 'finalizado' then
				update rendimiento_operadore set
	            duracion_atenciones_finalizadas = duracion_atenciones_finalizadas + duracion_llamada,
	            cantidad_atenciones_finalizadas = cantidad_atenciones_finalizadas + 1,
	            duracion_promedio_atenciones_finalizadas = (duracion_atenciones_finalizadas + duracion_llamada) / (cantidad_atenciones_finalizadas +1)
	       		where id_operadore = old.id_operadore and fecha_atencion = v_fecha_atencion;
			end if;
			
		end if;
	return new;
end;
$$ language plpgsql;

create or replace function enviar_email() returns trigger as $$
declare
    datos_cliente record;
    v_cuerpo text;
    v_asunto text;
begin
    select nombre, apellido, dni, telefono, email 
    into datos_cliente 
    from cliente
    where id_cliente = new.id_cliente;

	if (new.estado = 'iniciado') then
	    v_asunto := 'Skynet - nuevo trámite: ' || new.id_tramite;
	
	    v_cuerpo := 'CLIENTE: ' || datos_cliente.nombre || ' ' || datos_cliente.apellido || E'\n' ||
	                'DNI: ' || datos_cliente.dni || E'\n' ||
	                'TELÉFONO: ' || datos_cliente.telefono || E'\n' ||
	                'TRÁMITE: ' || new.tipo_tramite || E'\n' ||
	                'FECHA INICIO GESTIÓN: ' || new.f_inicio_gestion || E'\n' ||
	                'DESCRIPCIÓN: ' || new.descripcion;
	
	    insert into envio_email values (default, current_timestamp(0)::timestamp, datos_cliente.email,
	            v_asunto, v_cuerpo, null, null);
	
	else
			v_asunto := 'Skynet - cierre de trámite: ' || new.id_tramite;
		
		    v_cuerpo := 'CLIENTE: ' || datos_cliente.nombre || ' ' || datos_cliente.apellido || E'\n' ||
		                'DNI: ' || datos_cliente.dni || E'\n' ||
		                'TELÉFONO: ' || datos_cliente.telefono || E'\n' ||
		                'TRÁMITE: ' || new.tipo_tramite || E'\n' ||
		                'FECHA INICIO / FIN DE GESTIÓN: ' || new.f_inicio_gestion  ||' - ' || new.f_fin_gestion  || E'\n' ||
		                'DESCRIPCIÓN: ' || new.descripcion ||
						'ESTADO DE CIERRE: ' || new.estado || E'\n' ||
						'RESPUESTA SKYNET: ' || new.respuesta;
		
		    insert into envio_email values (default, current_timestamp(0)::timestamp, datos_cliente.email,
		            v_asunto, v_cuerpo, null, null);

	end if;
    return new;
end;
$$ language plpgsql;


--TRIGGERS

create or replace trigger reporte_rendimiento_trg
after update on cola_atencion
for each row
execute procedure reporte_rendimiento();

create or replace trigger enviar_email_trg
after insert or update on tramite
for each row
execute procedure enviar_email();
