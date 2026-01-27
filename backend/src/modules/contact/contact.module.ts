import { Module } from '@nestjs/common';
import { TypeOrmModule } from '@nestjs/typeorm';
import { ContactController } from './contact.controller';
import { ContactService } from './contact.service';
import { ConsultationRequest } from './entities/consultation-request.entity';

@Module({
  imports: [TypeOrmModule.forFeature([ConsultationRequest])],
  controllers: [ContactController],
  providers: [ContactService],
})
export class ContactModule {}
